package command

import (
	"net"

	"github.com/MagnusTiberius/packet"
)

// NewSubscribePacket creates a new SUBSCRIBE packet.
func NewSubscribePacket() *packet.SubscribePacket {
	return &packet.SubscribePacket{}
}

// NewUnsubscribePacket creates a new SUBSCRIBE packet.
func NewUnsubscribePacket() *packet.UnsubscribePacket {
	return &packet.UnsubscribePacket{}
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
	conn.Write([]byte("\n"))
}

//Unsubscribe send command to broker.
func Unsubscribe(packet *packet.UnsubscribePacket, conn net.Conn) {
	// Allocate buffer.
	buf := make([]byte, packet.Len())

	// Encode the packet.
	if _, err := packet.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	conn.Write(buf)
	conn.Write([]byte("\n"))
	conn.Write([]byte("\n"))
}
