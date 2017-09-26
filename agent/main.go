package main

import (
	"fmt"
	"net"

	"github.com/MagnusTiberius/meshnet/api/client"
	"github.com/gomqtt/packet"
)

func main() {

	cfg := client.ConfigTLS{
		Addr:      "127.0.0.1:8000",
		ClientPEM: "secure/certs/client.pem",
		ClientKey: "secure/certs/client.key",
	}

	c := client.NewClient()
	tls := c.NewTLS(&cfg)

	if tls == nil {
		return
	}

	// Connect
	p := packet.NewConnectPacket()
	p.Username = "gomqtt"
	p.Password = "amazing!"

	c.Connect(p, tls)

	c.HandleIncoming = handleIncoming

	fmt.Println("Calling HandleReceive")
	c.HandleReceive(tls)

}

func handleIncoming(buf []byte, conn net.Conn) {
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
