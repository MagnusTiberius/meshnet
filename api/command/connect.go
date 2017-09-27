package command

import (
	"net"

	"github.com/MagnusTiberius/packet"
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
