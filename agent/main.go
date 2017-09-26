package main

import (
	"fmt"
	"net"

	"github.com/MagnusTiberius/meshnet/api/client"
	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/gomqtt/packet"
)

var (
	uid uint16 = 1
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

	c.HandleIncoming = handleIncoming
	fmt.Println("Calling HandleReceive")
	go c.HandleReceive(tls)

	// Connect
	p := packet.NewConnectPacket()
	p.Username = "gomqtt"
	p.Password = "amazing!"

	fmt.Println("Calling Connect")
	c.Connect(p, tls)

	msg := command.NewMessage()
	msg.Topic = "welcome/all"
	msg.Payload = []byte("Hey this is a hello message.\n")

	fmt.Println("Calling Publish")
	command.Publish(msg, tls)

	msg.Payload = []byte("Another comment going in.\n")
	fmt.Println("Calling Publish")
	command.Publish(msg, tls)

	subp := command.NewSubscribePacket()
	subp.PacketID = uid
	uid = uid + 1
	sub := command.NewSubscription()
	sub.Topic = "welcome/all"
	subp.Subscriptions = append(subp.Subscriptions, sub)
	sub2 := command.NewSubscription()
	sub2.Topic = "goodbye/all"
	subp.Subscriptions = append(subp.Subscriptions, sub2)
	fmt.Printf("Calling Subscribe: %v \n", subp)
	command.Subscribe(subp, tls)

	msg.Payload = []byte("Another comment going in.\n")
	fmt.Println("Calling Publish")
	command.Publish(msg, tls)

	for {

	}

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
