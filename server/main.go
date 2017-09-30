package main

import (
	"io/ioutil"
	"log"
	"net"
	"os"

	"github.com/MagnusTiberius/meshnet/api/server"
	"github.com/MagnusTiberius/packet"
)

var (
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

	fh := server.FuncHandler{
		OnConnect:     OnConnect,
		OnPublish:     OnPublish,
		OnSubscribe:   OnSubscribe,
		OnDisconnect:  OnDisconnect,
		OnPingRequest: OnPingRequest,
	}

	go server.Start(broker, fh)

	for {
		go server.HandleConn(<-conns, broker)
	}

}

//OnConnect todo ...
func OnConnect(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnConnect\n")
}

//OnPublish todo ...
func OnPublish(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnPublish\n")
}

//OnSubscribe todo ...
func OnSubscribe(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnSubscribe\n")
}

//OnDisconnect todo ...
func OnDisconnect(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnDisconnect\n")
}

//OnPingRequest todo ...
func OnPingRequest(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnPingRequest\n")
}
