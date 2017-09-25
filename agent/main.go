package main

import (
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"

	"github.com/gomqtt/packet"
)

func main() {
	var n int
	cert, err := tls.LoadX509KeyPair("secure/certs/client.pem", "secure/certs/client.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", "127.0.0.1:8000", &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
	}
	defer conn.Close()
	log.Println("client: connected to: ", conn.RemoteAddr())

	state := conn.ConnectionState()
	for _, v := range state.PeerCertificates {
		fmt.Printf("PublicKey:\n")
		fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
		fmt.Printf("Subject:\n")
		fmt.Println(v.Subject)
	}
	log.Println("client: handshake: ", state.HandshakeComplete)
	log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)

	/*
		message := "Hello\n"
		n, err := io.WriteString(conn, message)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: wrote %q (%d bytes)", message, n)
	*/

	// Create new packet.
	pkt1 := packet.NewConnectPacket()
	pkt1.Username = "gomqtt"
	pkt1.Password = "amazing!"

	// Allocate buffer.
	buf := make([]byte, pkt1.Len())

	// Encode the packet.
	if _, err = pkt1.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	buf2 := fmt.Sprintf("%v\n", string(buf))
	conn.Write([]byte(buf2))

	reply := make([]byte, 4096)
	for {
		n, err = conn.Read(reply)
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
		handleIncomin(reply, conn)
	}

	//log.Print("client: exiting")
}

func handleIncomin(buf []byte, conn net.Conn) {
	// Detect packet.
	l, mt := packet.DetectPacket(buf)

	// Check length
	if l == 0 {
		fmt.Printf("buffer not complete yet")
		return // buffer not complete yet
	}

	// Create packet.
	pkt2, err := mt.New()
	if err != nil {
		panic(err) // packet type is invalid
	}

	// Decode packet.
	_, err = pkt2.Decode(buf)
	if err != nil {
		fmt.Printf("%v\n", err)
		panic(err) // there was an error while decoding
	}

	switch pkt2.Type() {
	case packet.CONNECT:
		c := pkt2.(*packet.ConnectPacket)
		fmt.Println(c.Username)
		fmt.Println(c.Password)
	case packet.CONNACK:
		ack := pkt2.(*packet.ConnackPacket)
		fmt.Printf("ReturnCode:%v\n\n", ack.ReturnCode)
	}
}
