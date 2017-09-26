package main

import (
	"fmt"
	"log"
	"net"

	"github.com/MagnusTiberius/meshnet/api"
	"github.com/MagnusTiberius/meshnet/api/client"
	"github.com/gomqtt/packet"
)

func main() {

	cfg := client.Config{
		Addr:      "127.0.0.1:8000",
		ClientPEM: "secure/certs/client.pem",
		ClientKey: "secure/certs/client.key",
	}
	client := client.NewClientTLS(&cfg)

	if client == nil {
		return
	}

	// Connect
	packet := packet.NewConnectPacket()
	packet.Username = "gomqtt"
	packet.Password = "amazing!"

	command.Connect(packet, client.Conn)

	go handleReceive(client.Conn)

	handleReceive(client.Conn)

	//log.Print("client: exiting")
}

func initSenders(conn net.Conn) {

}

func handleReceive(conn net.Conn) {
	reply := make([]byte, 4096)
	for {
		n, err := conn.Read(reply)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
		handleIncomin(reply, conn)
	}
}

func handleIncomin(buf []byte, conn net.Conn) {
	// Detect packet.
	l, mt := packet.DetectPacket(buf)

	// Check length
	if l == 0 {
		fmt.Printf("buffer not complete yet")
		return // buffer not complete yet
	}

	// Create packet.
	pkt2, err := mt.New()
	if err != nil {
		panic(err) // packet type is invalid
	}

	// Decode packet.
	_, err = pkt2.Decode(buf)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err) // there was an error while decoding
	}

	switch pkt2.Type() {
	case packet.CONNECT:
		c := pkt2.(*packet.ConnectPacket)
		fmt.Println(c.Username)
		fmt.Println(c.Password)
	case packet.CONNACK:
		ack := pkt2.(*packet.ConnackPacket)
		fmt.Printf("ReturnCode:%v\n\n", ack.ReturnCode)
	}
}
