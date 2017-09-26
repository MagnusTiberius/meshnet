package main

import (
	"net"
	"os"

	"github.com/MagnusTiberius/meshnet/api/server"
)

var (
	connPool  map[string]*net.Conn
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
	broker.Accept()

}
