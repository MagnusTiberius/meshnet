package server

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/MagnusTiberius/packet"
)

var (
	connPool    map[string]net.Conn
	port        string
	ip          string
	relayConn   net.Conn
	tlsOk       string
	funcHandler FuncHandler
)

//Config todo ...
type Config struct {
	Addr string
	PEM  string //"secure/certs/client.pem"
	Key  string //"secure/certs/client.key"
}

//FuncHandler todo ...
type FuncHandler struct {
	OnConnect     func(conn net.Conn, pkt packet.Packet)
	OnPublish     func(conn net.Conn, pkt packet.Packet)
	OnSubscribe   func(conn net.Conn, pkt packet.Packet)
	OnDisconnect  func(conn net.Conn, pkt packet.Packet)
	OnPingRequest func(conn net.Conn, pkt packet.Packet)
}

//ListenerTLS todo ...
func ListenerTLS(cfg *Config) net.Listener {
	cert, err := tls.LoadX509KeyPair(cfg.PEM, cfg.Key)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
		return nil
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", cfg.Addr, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
		return nil
	}
	fmt.Printf("Server listening...\n")
	return listener
}

//HandleConns todo ...
func HandleConns(l net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			cl, err := l.Accept()
			if cl == nil {
				panic(fmt.Sprintf("%s: %v", "Listener Accept() failure: ", err))
				//continue
			}
			i++
			//log.Printf("%d: %v accepted %v\n", i, cl.LocalAddr(), cl.RemoteAddr())
			//conn_pool[fmt.Sprintf("%v", client.RemoteAddr())] = &client
			//cl.Write([]byte("Welcome to echoserver utopia\n"))
			//handleEvent(Event{Name: "CONNECT_EVENT", Client: client})
			ch <- cl
		}
	}()
	return ch
}

//Start todo ...
func Start(b *Broker, fh FuncHandler) {
	funcHandler = fh
	ctr := 0
	for {
		ctr = ctr + 1
		time.Sleep(1000 * time.Millisecond)
		//Walk the connection pool and check each network connection.
		if ctr > 5 {
			//five seconds has elapsed
			for kcn, kv := range connPool {
				if kv != nil {
					_, err := pingClient(kv)
					if err != nil {
						//Connection is lost, remove it from the pool/list.
						log.Printf("Conn Closed: %v \n", kcn)
						for _, v := range b.Bundle.TopicList {
							log.Printf("Removing element %v\n", kcn)
							delete(v.ConnList, kcn)
						}
						connPool[kcn] = nil
					}
				}
			}
			ctr = 0
		}
		//Now, dispatch the messages to the subscribers
		for key, v := range b.Bundle.TopicList {
			log.Printf("key: %v \n", key)
			for _, d := range v.ConnList {
				addr := d.RemoteAddr()
				log.Printf("\taddr: %v \n", addr)
				h := v.LastSent
				g := len(v.Messages)
				for k, m := range v.Messages {
					if int64(k) > h {
						log.Printf("\t\tmsg: %v \n", string(m.Payload))
						_, err := command.Publish(m, d)
						if err != nil {
							log.Printf("\t\t\t Invalid Address\n")
						}
					}
				}
				v.LastSent = int64(g)
			}
		}
	}
}

//HandleConn  will handle the connection from the channel
func HandleConn(c net.Conn, broker *Broker) {
	b := bufio.NewReader(c)
	for {
		msg, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		handleIncoming(msg, c, broker)
	}
}

//handleIncoming todo ...
func handleIncoming(buf []byte, conn net.Conn, brk *Broker) {
	// Detect packet.
	l, mt := packet.DetectPacket(buf)

	// Check length
	if l == 0 {
		log.Printf("buffer not complete yet")
		return // buffer not complete yet
	}

	// Create packet.
	pkt, err := mt.New()
	if err != nil {
		panic(err) // packet type is invalid
	}

	// Decode packet.
	_, err = pkt.Decode(buf)
	if err != nil {
		panic(err) // there was an error while decoding
	}

	switch pkt.Type() {
	case packet.PINGREQ:
		pingReply(conn)
		if funcHandler.OnPingRequest != nil {
			funcHandler.OnPingRequest(conn, pkt)
		}
	case packet.PINGRESP:
		log.Printf("Ping response\n")
	case packet.CONNECT:
		log.Printf("\nCONNECT:\n")
		c := pkt.(*packet.ConnectPacket)
		log.Println("Username:" + c.Username)
		log.Println("Password:" + c.Password)
		addr := fmt.Sprintf("%v", conn.RemoteAddr())
		if connPool == nil {
			connPool = make(map[string]net.Conn)
		}
		connPool[addr] = conn
		replyConnectionAck(conn)
		if funcHandler.OnConnect != nil {
			funcHandler.OnConnect(conn, pkt)
		}
	case packet.PUBLISH:
		log.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		log.Println("Topic:" + p.Message.Topic)
		log.Println("Payload:" + string(p.Message.Payload))
		brk.Bundle.Publish(&p.Message, conn)
		if funcHandler.OnPublish != nil {
			funcHandler.OnPublish(conn, pkt)
		}
	case packet.SUBSCRIBE:
		log.Printf("\nSUBSCRIBE:\n")
		p := pkt.(*packet.SubscribePacket)
		log.Printf("Subscriptions: %v \n", p.Subscriptions)
		replySubscriptionAck(conn, p.PacketID)
		for _, s := range p.Subscriptions {
			brk.Bundle.Subscribe(&s, conn)
		}
		if funcHandler.OnSubscribe != nil {
			funcHandler.OnSubscribe(conn, pkt)
		}
	case packet.DISCONNECT:
		conn.Close()
		if funcHandler.OnDisconnect != nil {
			funcHandler.OnDisconnect(conn, pkt)
		}
	}

}

func replySubscriptionAck(c net.Conn, uid uint16) {
	log.Println("replySubscriptionAck")
	ack := packet.NewSubackPacket()
	ack.PacketID = uid
	ack.ReturnCodes = []uint8{0}
	//ack.ReturnCode = []byte{0, 1}
	//ack.SessionPresent = true

	// Allocate buffer.
	buf := make([]byte, ack.Len())

	// Encode the packet.
	if _, err := ack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	c.Write(buf)
	c.Write([]byte("\n"))
}

//pingClient todo ...
func pingClient(c net.Conn) (n int, err error) {
	log.Printf("PingClient %v\n", c.RemoteAddr())
	ping := packet.NewPingreqPacket()

	// Allocate buffer.
	buf := make([]byte, ping.Len())

	// Encode the packet.
	if _, err = ping.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	return n, err
}

//pingReply todo ...
func pingReply(c net.Conn) (n int, err error) {
	log.Println("pingReply")
	pingack := packet.NewPingrespPacket()

	// Allocate buffer.
	buf := make([]byte, pingack.Len())

	// Encode the packet.
	if _, err = pingack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	return n, err
}

func replyConnectionAck(c net.Conn) (n int, err error) {
	log.Println("replyConnectionAck")
	ack := packet.NewConnackPacket()
	ack.ReturnCode = packet.ConnectionAccepted
	ack.SessionPresent = true

	// Allocate buffer.
	buf := make([]byte, ack.Len())

	// Encode the packet.
	if _, err = ack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	log.Printf("replyConnectionAck...done: %v \n", buf)
	return n, err
}
