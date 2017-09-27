package command

import (
	"net"

	"github.com/MagnusTiberius/packet"
)

//Publish send public to broker
func Publish(message *packet.Message, conn net.Conn) (n int, err error) {
	pub := packet.NewPublishPacket()
	pub.Message = *message
	pub.Dup = false

	return publishPacket(pub, conn)
}

//publishPacket send public to broker
func publishPacket(packet *packet.PublishPacket, conn net.Conn) (n int, err error) {
	// Allocate buffer.
	buf := make([]byte, packet.Len())

	// Encode the packet.
	if _, err = packet.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	n, err = conn.Write(buf)
	conn.Write([]byte("\n"))
	return n, err
}
