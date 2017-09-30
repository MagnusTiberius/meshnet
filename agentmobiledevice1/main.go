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

const (
	ulimit = 10
)

var (
	uid  uint16 = 100
	ctr  int
	once bool
	tls  net.Conn
)

func main() {
	log.SetOutput(ioutil.Discard)

	cfg := client.ConfigTLS{
		Addr:      "127.0.0.1:8000",
		ClientPEM: "secure/certs/client.pem",
		ClientKey: "secure/certs/client.key",
	}
	tls = client.NewTLS(&cfg)

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
		OnUnsubscribe: OnUnsubscribe,
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
	ctr = ctr + 1
	if ctr > ulimit {
		if !once {
			u := command.NewUnsubscribePacket()
			u.PacketID = 100
			u.Topics = []string{"device/sensor/2"}
			command.Unsubscribe(u, tls)
			once = !once
		}
	}
}

//OnSubscribe todo ...
func OnSubscribe(conn net.Conn, pkt packet.Packet) {
	//log.Printf("OnSubscribe\n")
}

//OnUnsubscribe todo ...
func OnUnsubscribe(conn net.Conn, pkt packet.Packet) {
	fmt.Printf("OnUnsubscribe\n")
}

//OnPingRequest todo ...
func OnPingRequest(conn net.Conn, pkt packet.Packet) {
	fmt.Printf("OnPingRequest\n")
}
