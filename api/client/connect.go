package client

import (
	"bufio"
	"crypto/tls"
	"log"
	"net"
	"time"

	"github.com/MagnusTiberius/packet"
)

//ConfigTLS todo ...
type ConfigTLS struct {
	Addr      string
	ClientPEM string //"secure/certs/client.pem"
	ClientKey string //"secure/certs/client.key"
}

//NewTLS todo...
func NewTLS(cfg *ConfigTLS) net.Conn {
	cert, err := tls.LoadX509KeyPair(cfg.ClientPEM, cfg.ClientKey)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
		panic(err)
		//return nil
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", cfg.Addr, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
		panic(err)
		//return nil
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	/*
		state := conn.ConnectionState()
		for _, v := range state.PeerCertificates {
			fmt.Printf("PublicKey:\n")
			fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
			fmt.Printf("Subject:\n")
			fmt.Println(v.Subject)
		}
		log.Println("client: handshake: ", state.HandshakeComplete)
		log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	*/
	return conn
}

//HandleReceive todo ...
func HandleReceive(conn net.Conn) {
	br := bufio.NewReader(conn)
	for {
		msg, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("-")
		HandleIncoming(msg, conn)
	}
}

//HandleIncoming todo ...
func HandleIncoming(buf []byte, conn net.Conn) {
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
		c := pkt.(*packet.ConnectPacket)
		log.Println(c.Username)
		log.Println(c.Password)
	case packet.CONNACK:
		ack := pkt.(*packet.ConnackPacket)
		log.Printf("ReturnCode:%v\n\n", ack.ReturnCode)
	case packet.SUBACK:
		log.Printf("SUBACK:\n\n")
	case packet.PUBLISH:
		log.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		log.Printf("%v, Topic: %v, Payload: %v \n ", time.Now(), p.Message.Topic, string(p.Message.Payload))
	}
}

func disconnect(c net.Conn) (n int, err error) {
	log.Println("NewDisconnectPacket")
	discon := packet.NewDisconnectPacket()

	// Allocate buffer.
	buf := make([]byte, discon.Len())

	// Encode the packet.
	if _, err = discon.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	return n, err
}

//pingReply todo ...
func pingReply(c net.Conn) (n int, err error) {
	log.Println("PingReply")
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
