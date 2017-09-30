package client

import (
	"bufio"
	"log"
	"net"

	"github.com/MagnusTiberius/packet"
)

//HandleReceive todo ...
func HandleReceive(conn net.Conn, fh FuncHandler) {
	funcHandler = fh
	br := bufio.NewReader(conn)
	for {
		msg, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("-")
		HandleIncoming(msg, conn)
	}
}

//HandleIncoming todo ...
func HandleIncoming(buf []byte, conn net.Conn) {
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
		if funcHandler.OnPingRequest != nil {
			funcHandler.OnPingRequest(conn, pkt)
		}
	case packet.PINGRESP:
		log.Printf("Ping response\n")
	case packet.CONNECT:
		c := pkt.(*packet.ConnectPacket)
		log.Println(c.Username)
		log.Println(c.Password)
		if funcHandler.OnConnect != nil {
			funcHandler.OnConnect(conn, pkt)
		}
	case packet.CONNACK:
		//ack := pkt.(*packet.ConnackPacket)
	case packet.SUBACK:
	case packet.UNSUBACK:
		if funcHandler.OnUnsubscribe != nil {
			funcHandler.OnUnsubscribe(conn, pkt)
		}
	case packet.PUBLISH:
		//p := pkt.(*packet.PublishPacket)
		if funcHandler.OnPublish != nil {
			funcHandler.OnPublish(conn, pkt)
		}
	}
}

//Disconnect todo ...
func Disconnect(c net.Conn) (n int, err error) {
	log.Println("disconnect")
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
