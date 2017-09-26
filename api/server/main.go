package server

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"github.com/gomqtt/packet"
)

//Broker todo ...
type Broker struct {
	Listener net.Listener
}

//NewBroker todo ...
func NewBroker() *Broker {
	return &Broker{}
}

//Accept todo ...
func (b *Broker) Accept() {
	for {
		conn, err := b.Listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			fmt.Printf("%v\n", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		go b.HandleConn(conn)
	}
}

//HandleConn todo ...
func (b *Broker) HandleConn(conn net.Conn) {
	br := bufio.NewReader(conn)
	for {
		msg, err := br.ReadBytes('\n')
		if err != nil {
			break
		}

		b.HandleIncoming(msg, conn)
	}
}

//HandleIncoming todo ...
func (b *Broker) HandleIncoming(buf []byte, conn net.Conn) {
	// Detect packet.
	l, mt := packet.DetectPacket(buf)

	// Check length
	if l == 0 {
		fmt.Printf("buffer not complete yet")
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
	case packet.CONNECT:
		c := pkt.(*packet.ConnectPacket)
		fmt.Println(c.Username)
		fmt.Println(c.Password)
		replyConnectionAck(conn)
	}

}

func replyConnectionAck(c net.Conn) {
	ack := packet.NewConnackPacket()
	ack.ReturnCode = packet.ConnectionAccepted
	ack.SessionPresent = true

	// Allocate buffer.
	buf := make([]byte, ack.Len())

	// Encode the packet.
	if _, err := ack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	c.Write(buf)
}
