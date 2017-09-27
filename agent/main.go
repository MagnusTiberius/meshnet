package main

import (
	"bufio"
	"fmt"
	"log"
	"net"
	"time"

	"github.com/MagnusTiberius/meshnet/api/client"
	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/gomqtt/packet"
)

var (
	uid uint16 = 100
)

func main() {
	fmt.Printf("00001\n")
	cfg := client.ConfigTLS{
		Addr:      "127.0.0.1:8000",
		ClientPEM: "secure/certs/client.pem",
		ClientKey: "secure/certs/client.key",
	}
	tls := client.NewTLS(&cfg)
	fmt.Printf("00002\n")

	if tls == nil {
		panic("null connection")
	}

	// Connect
	packet := packet.NewConnectPacket()
	packet.Username = "gomqtt"
	packet.Password = "amazing!"

	command.Connect(packet, tls)

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
	sub3 := command.NewSubscription()
	sub3.Topic = "awesome/all"
	subp.Subscriptions = append(subp.Subscriptions, sub3)
	fmt.Printf("Calling Subscribe: %v \n", subp)
	command.Subscribe(subp, tls)

	msg.Payload = []byte("I would like to see all of them.\n")
	fmt.Println("Calling Publish")
	command.Publish(msg, tls)

	go handleReceive(tls)

	handleReceive(tls)

	//log.Print("client: exiting")
}

func initSenders(conn net.Conn) {

}

func handleReceive(conn net.Conn) {
	//reply := make([]byte, 4096)
	br := bufio.NewReader(conn)
	for {
		//n, err := conn.Read(reply)
		msg, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		//log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
		handleIncomin(msg, conn)
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
	pkt, err := mt.New()
	if err != nil {
		panic(err) // packet type is invalid
	}

	//fmt.Printf("BUF: %v, %v \n", buf, string(buf))

	// Decode packet.
	_, err = pkt.Decode(buf)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err) // there was an error while decoding
	}

	switch pkt.Type() {
	case packet.CONNECT:
		c := pkt.(*packet.ConnectPacket)
		fmt.Println(c.Username)
		fmt.Println(c.Password)
	case packet.CONNACK:
		ack := pkt.(*packet.ConnackPacket)
		fmt.Printf("ReturnCode:%v\n\n", ack.ReturnCode)
	case packet.SUBACK:
		fmt.Printf("SUBACK:\n\n")
		sub := pkt.(*packet.SubackPacket)
		fmt.Printf("PacketID:%v\n\n", sub.PacketID)
		fmt.Printf("ReturnCodes:%v\n\n", sub.ReturnCodes)
	case packet.PUBLISH:
		fmt.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		fmt.Printf("%v, Topic: %v, Payload: %v \n ", time.Now(), p.Message.Topic, string(p.Message.Payload))
		//fmt.Println("\tPayload:" + string(p.Message.Payload))
		//brk.Bundle.Publish(&p.Message, conn)
	}
}
