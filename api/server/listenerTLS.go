package server

import (
	"crypto/rand"
	"crypto/tls"
	"fmt"
	"log"
	"net"
)

//Config todo ...
type Config struct {
	Addr string
	PEM  string //"secure/certs/client.pem"
	Key  string //"secure/certs/client.key"
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
			client, err := l.Accept()
			if client == nil {
				panic(fmt.Sprintf("%s: %v", "Listener Accept() failure: ", err))
				//continue
			}
			i++
			fmt.Printf("%d: %v accepted %v\n", i, client.LocalAddr(), client.RemoteAddr())
			//conn_pool[fmt.Sprintf("%v", client.RemoteAddr())] = &client
			client.Write([]byte("Welcome to echoserver utopia\n"))
			//handleEvent(Event{Name: "CONNECT_EVENT", Client: client})
			ch <- client
		}
	}()
	return ch
}
