package server

import (
	"bufio"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"github.com/MagnusTiberius/meshnet/api/repo"
)

//Broker todo ...
type Broker struct {
	Listener       net.Listener
	HandleIncoming func(buf []byte, conn net.Conn, brk *Broker)
	Bundle         *repo.Bundle
}

//NewBroker todo ...
func NewBroker() *Broker {
	return &Broker{
		Bundle: repo.NewBundle(),
	}
}

//HandleConns todo ...
func (b *Broker) HandleConns(l net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := l.Accept()
			if client == nil {
				panic(fmt.Sprintf("%s: %v", "Listener Accept() failure: ", err))
				//continue
			}
			i++
			fmt.Printf("%d: %v accepted %v\n", i, client.LocalAddr(), client.RemoteAddr())
			//conn_pool[fmt.Sprintf("%v", client.RemoteAddr())] = &client
			//client.Write([]byte("Welcome to echoserver utopia\n"))
			//handleEvent(Event{Name: "CONNECT_EVENT", Client: client})
			ch <- client
		}
	}()
	return ch
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

		go b.HandleIncoming(msg, conn, b)
	}
}
