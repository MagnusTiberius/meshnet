package command

import (
	"net"

	"github.com/gomqtt/packet"
)

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
