package server

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
)

//Broker todo ...
type Broker struct {
	Listener       net.Listener
	HandleIncoming func(buf []byte, conn net.Conn)
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
		b.HandleConn(conn)
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

		go b.HandleIncoming(msg, conn)
	}
}
