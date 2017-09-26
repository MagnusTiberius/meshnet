package main

import (
	"fmt"
	"net"
	"os"

	"github.com/MagnusTiberius/meshnet/api/server"
	"github.com/gomqtt/packet"
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
	broker.HandleIncoming = handleIncoming
	broker.Accept()

}

//handleIncoming todo ...
func handleIncoming(buf []byte, conn net.Conn) {
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
		fmt.Printf("\nCONNECT:\n")
		c := pkt.(*packet.ConnectPacket)
		fmt.Println("Username:" + c.Username)
		fmt.Println("Password:" + c.Password)
		replyConnectionAck(conn)
	case packet.PUBLISH:
		fmt.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		fmt.Println("Topic:" + p.Message.Topic)
		fmt.Println("Payload:" + string(p.Message.Payload))
	case packet.SUBSCRIBE:
		fmt.Printf("\nSUBSCRIBE:\n")
		p := pkt.(*packet.SubscribePacket)
		fmt.Printf("Subscriptions: %v \n", p.Subscriptions)
	}

}

func replyConnectionAck(c net.Conn) {
	fmt.Println("replyConnectionAck")
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
	c.Write([]byte("\n"))
	fmt.Printf("replyConnectionAck...done: %v \n", buf)
}
