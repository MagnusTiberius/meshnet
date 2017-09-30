package main

import (
	"fmt"
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

	subp := command.NewSubscribePacket()
	subp.PacketID = uid
	uid = uid + 1

	sub := command.NewSubscription()
	sub.Topic = "device/sensor/1"
	subp.Subscriptions = append(subp.Subscriptions, sub)

	sub = command.NewSubscription()
	sub.Topic = "device/sensor/2"
	subp.Subscriptions = append(subp.Subscriptions, sub)

	log.Printf("Calling Subscribe: %v \n", subp)
	command.Subscribe(subp, tls)

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
	//log.Printf("OnConnect\n")
}

//OnPublish todo ...
func OnPublish(conn net.Conn, pkt packet.Packet) {
	p := pkt.(*packet.PublishPacket)
	fmt.Printf("Topic:%v; Payload:%v\n", p.Message.Topic, string(p.Message.Payload))
}

//OnSubscribe todo ...
func OnSubscribe(conn net.Conn, pkt packet.Packet) {
	//log.Printf("OnSubscribe\n")
}

//OnPingRequest todo ...
func OnPingRequest(conn net.Conn, pkt packet.Packet) {
	//log.Printf("OnPingRequest\n")
}
