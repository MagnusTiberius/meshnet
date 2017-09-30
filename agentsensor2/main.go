package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"math/rand"
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

	fh := client.FuncHandler{
		OnConnect:     OnConnect,
		OnPublish:     OnPublish,
		OnSubscribe:   OnSubscribe,
		OnPingRequest: OnPingRequest,
	}

	go client.HandleReceive(tls, fh)

	//client.HandleReceive(tls, fh)

	msg1 := command.NewMessage()
	msg1.Topic = "device/sensor/2"
	for {
		num := rand.Float64() * 40
		msg1.Payload = []byte(fmt.Sprintf("{\"sensor\":\"water flow meter\", \"value\":%v}\n", num))
		command.Publish(msg1, tls)
		time.Sleep(1200 * time.Millisecond)
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

//OnPingRequest todo ...
func OnPingRequest(conn net.Conn, pkt packet.Packet) {
	log.Printf("OnPingRequest\n")
}
