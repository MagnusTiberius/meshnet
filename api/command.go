package command

import (
	"net"

	"github.com/gomqtt/packet"
)

//Connect send public to broker
func Connect(packet *packet.ConnectPacket, conn net.Conn) {
	// Allocate buffer.
	buf := make([]byte, packet.Len())

	// Encode the packet.
	if _, err := packet.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	conn.Write(buf)
	conn.Write([]byte("\n"))
}

//NewMessage returns a new packet message
func NewMessage() *packet.Message {
	return &packet.Message{}
}

//Publish send public to broker
func Publish(message *packet.Message, conn net.Conn) {
	pub := packet.NewPublishPacket()
	pub.Message = *message
	pub.Dup = false

	publishPacket(pub, conn)
}

//publishPacket send public to broker
func publishPacket(packet *packet.PublishPacket, conn net.Conn) {
	// Allocate buffer.
	buf := make([]byte, packet.Len())

	// Encode the packet.
	if _, err := packet.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	conn.Write(buf)
	conn.Write([]byte("\n"))
}

// NewSubscribePacket creates a new SUBSCRIBE packet.
func NewSubscribePacket() *packet.SubscribePacket {
	return &packet.SubscribePacket{}
}

//NewSubscription returns a new Sub
func NewSubscription() packet.Subscription {
	return packet.Subscription{}
}

//Subscribe send command to broker.
func Subscribe(packet *packet.SubscribePacket, conn net.Conn) {
	// Allocate buffer.
	buf := make([]byte, packet.Len())

	// Encode the packet.
	if _, err := packet.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	conn.Write(buf)
	conn.Write([]byte("\n"))
}
