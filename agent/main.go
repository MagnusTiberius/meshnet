package main

import (
	"bufio"
	"io/ioutil"
	"log"
	"net"
	"time"

	"github.com/MagnusTiberius/meshnet/api/client"
	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/MagnusTiberius/packet"
)

var (
	uid uint16 = 100
)

func main() {
	log.SetOutput(ioutil.Discard)

	cfg := client.ConfigTLS{
		Addr:      "127.0.0.1:8000",
		ClientPEM: "secure/certs/client.pem",
		ClientKey: "secure/certs/client.key",
	}
	tls := client.NewTLS(&cfg)

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

	log.Println("Calling Publish")
	command.Publish(msg, tls)

	msg.Payload = []byte("Another comment going in.\n")
	log.Println("Calling Publish")
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
	log.Printf("Calling Subscribe: %v \n", subp)
	command.Subscribe(subp, tls)

	msg.Payload = []byte("I would like to see all of them.\n")
	log.Println("Calling Publish")
	command.Publish(msg, tls)

	go handleReceive(tls)

	handleReceive(tls)

	//log.Print("client: exiting")
}

func initSenders(conn net.Conn) {

}

func handleReceive(conn net.Conn) {
	br := bufio.NewReader(conn)
	for {
		msg, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("-")
		handleIncoming(msg, conn)
	}
}

func handleIncoming(buf []byte, conn net.Conn) {
	// Detect packet.
	l, mt := packet.DetectPacket(buf)

	// Check length
	if l == 0 {
		log.Printf("buffer not complete yet")
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
	case packet.PINGREQ:
		pingReply(conn)
	case packet.PINGRESP:
		log.Printf("Ping response\n")
	case packet.CONNECT:
		c := pkt.(*packet.ConnectPacket)
		log.Println(c.Username)
		log.Println(c.Password)
	case packet.CONNACK:
		ack := pkt.(*packet.ConnackPacket)
		log.Printf("ReturnCode:%v\n\n", ack.ReturnCode)
	case packet.SUBACK:
		log.Printf("SUBACK:\n\n")
	case packet.PUBLISH:
		log.Printf("\nPUBLISH:\n")
		p := pkt.(*packet.PublishPacket)
		log.Printf("%v, Topic: %v, Payload: %v \n ", time.Now(), p.Message.Topic, string(p.Message.Payload))
	}
}

func disconnect(c net.Conn) (n int, err error) {
	log.Println("NewDisconnectPacket")
	discon := packet.NewDisconnectPacket()

	// Allocate buffer.
	buf := make([]byte, discon.Len())

	// Encode the packet.
	if _, err = discon.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	return n, err
}

//pingReply todo ...
func pingReply(c net.Conn) (n int, err error) {
	log.Println("PingReply")
	pingack := packet.NewPingrespPacket()

	// Allocate buffer.
	buf := make([]byte, pingack.Len())

	// Encode the packet.
	if _, err = pingack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = c.Write(buf)
	c.Write([]byte("\n"))
	return n, err
}
