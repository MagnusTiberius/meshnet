package main

import (
	"bufio"
	"fmt"
	"io/ioutil"
	"log"
	"net"
	"os"
	"time"

	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/MagnusTiberius/meshnet/api/server"
	"github.com/MagnusTiberius/packet"
)

var (
	connPool  map[string]net.Conn
	port      string
	ip        string
	relayConn net.Conn
	tlsOk     string
)

//Event event type
type Event struct {
	Name   string
	Client net.Conn
	Msg    []byte
}

func main() {
	log.SetOutput(ioutil.Discard)

	if len(os.Args) == 2 {
		tlsOk = os.Args[1]
	}

	broker := server.NewBroker()

	cfg := server.Config{
		Addr: "127.0.0.1:8000",
		PEM:  "secure/certs/server.pem",
		Key:  "secure/certs/server.key",
	}
	broker.Listener = server.ListenerTLS(&cfg)

	conns := broker.HandleConns(broker.Listener)

	go startServer(broker)

	for {
		go handleConn(<-conns, broker)
	}

}

func handleConn(c net.Conn, broker *server.Broker) {
	b := bufio.NewReader(c)
	for {
		msg, err := b.ReadBytes('\n')
		if err != nil {
			break
		}
		handleIncoming(msg, c, broker)
	}
}

func startServer(b *server.Broker) {
	for {
		time.Sleep(1000 * time.Millisecond)
		log.Printf(".")
		for kcn, kv := range connPool {
			if kv != nil {
				_, err := pingClient(kv)
				if err != nil {
					log.Printf("Conn Closed: %v \n", kcn)
					for _, v := range b.Bundle.TopicList {
						log.Printf("Removing element %v\n", kcn)
						delete(v.ConnList, kcn)
					}
					connPool[kcn] = nil
				}
			}
		}
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

//handleIncoming todo ...
func handleIncoming(buf []byte, conn net.Conn, brk *server.Broker) {
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
	case packet.PINGRESP:
		log.Printf("Ping response\n")
	case packet.CONNECT:
		log.Printf("\nCONNECT:\n")
		c := pkt.(*packet.ConnectPacket)
		log.Println("Username:" + c.Username)
		log.Println("Password:" + c.Password)
		replyConnectionAck(conn)
		addr := fmt.Sprintf("%v", conn.RemoteAddr())
		if connPool == nil {
			connPool = make(map[string]net.Conn)
		}
		connPool[addr] = conn
	case packet.PUBLISH:
		log.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		log.Println("Topic:" + p.Message.Topic)
		log.Println("Payload:" + string(p.Message.Payload))
		brk.Bundle.Publish(&p.Message, conn)
	case packet.SUBSCRIBE:
		log.Printf("\nSUBSCRIBE:\n")
		p := pkt.(*packet.SubscribePacket)
		log.Printf("Subscriptions: %v \n", p.Subscriptions)
		replySubscriptionAck(conn, p.PacketID)
		for _, s := range p.Subscriptions {
			brk.Bundle.Subscribe(&s, conn)
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
