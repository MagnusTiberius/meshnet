package main

import (
	"io/ioutil"
	"log"
	"net"

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

	fh := client.FuncHandler{
		OnConnect:     OnConnect,
		OnPublish:     OnPublish,
		OnSubscribe:   OnSubscribe,
		OnPingRequest: OnPingRequest,
	}

	//go client.HandleReceive(tls, fh)

	client.HandleReceive(tls, fh)

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

//OnPingRequest todo ...
func OnPingRequest(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnPingRequest\n")
}
