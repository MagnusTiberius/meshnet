package main

import (
	"bufio"
	"crypto/rand"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"log"
	"net"
	"os"
	"strconv"

	"github.com/MagnusTiberius/meshnet/api/lex"
	"github.com/gomqtt/packet"
)

var (
	connPool  map[string]*net.Conn
	port      string
	ip        string
	relayConn net.Conn
	tlsOk     string
)

//Event event type
type Event struct {
	Name   string
	Client net.Conn
	Msg    []byte
}

func main() {
	if len(os.Args) == 2 {
		tlsOk = os.Args[1]
	}

	if tlsOk == "tls" {
		setupListenerTLS("127.0.0.1", "8000")
	} else {
		setupListener("127.0.0.1", "8049")
	}
}

func setupListener(ip string, sport string) {
	nport, _ := strconv.Atoi(sport)
	n := nport + 1

	addr := fmt.Sprintf("%s:%d", ip, n)
	server, err := net.Listen("tcp", addr)

	if server == nil {
		panic(fmt.Sprintf("%s: %v", "Listen failure: ", err))
	}
	connPool = map[string]*net.Conn{}
	conns := handleConns(server)
	for {
		go handleConn(<-conns)
	}

}

func handleConns(l net.Listener) chan net.Conn {
	ch := make(chan net.Conn)
	i := 0
	go func() {
		for {
			client, err := l.Accept()
			if client == nil {
				panic(fmt.Sprintf("%s: %v", "Listener Accept() failure: ", err))
				//continue
			}
			i++
			fmt.Printf("%d: %v accepted %v\n", i, client.LocalAddr(), client.RemoteAddr())
			connPool[fmt.Sprintf("%v", client.RemoteAddr())] = &client
			client.Write([]byte("Welcome to echoserver utopia\n"))
			handleEvent(Event{Name: "CONNECT_EVENT", Client: client})
			ch <- client
		}
	}()
	return ch
}

func handleEvent(e Event) {
	//fmt.Println(e.Name)
	switch e.Name {
	case "CONNECT_EVENT":
		for _, c := range connPool {
			conn := *c
			addr := fmt.Sprintf("%v", e.Client.RemoteAddr())
			if fmt.Sprintf("%v", conn.RemoteAddr()) != addr {
				msgConnect := []byte(fmt.Sprintf("%s has connected\n", addr))
				conn.Write(msgConnect)
			}
		}
	case "CLIENT_MSG":
		for _, c := range connPool {
			addr := fmt.Sprintf("%v", e.Client.RemoteAddr())
			conn := *c
			if fmt.Sprintf("%v", conn.RemoteAddr()) != addr {
				msg := fmt.Sprintf("ECHO>>%v:%s", e.Client.RemoteAddr(), string(e.Msg))
				conn.Write([]byte(msg))
			}
		}
	case "RELAY_MSG":
		for _, c := range connPool {
			client := *c
			msg := fmt.Sprintf("RELAY>>%v:%s\n", e.Client.RemoteAddr(), string(e.Msg))
			client.Write([]byte(msg))
		}
	default:
	}
	//relayEvent(e)
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
		panic(err) // there was an error while decoding
	}

	switch pkt2.Type() {
	case packet.CONNECT:
		c := pkt2.(*packet.ConnectPacket)
		fmt.Println(c.Username)
		fmt.Println(c.Password)
		handleConnectionAck(conn)
	}

}

func handleConnectionAck(c net.Conn) {
	ack := packet.NewConnackPacket()
	ack.ReturnCode = packet.ConnectionAccepted
	ack.SessionPresent = true

	// Allocate buffer.
	buf := make([]byte, ack.Len())

	// Encode the packet.
	if _, err := ack.Encode(buf); err != nil {
		panic(err) // error while encoding
	}

	c.Write(buf)
}

func handleConn(c net.Conn) {
	b := bufio.NewReader(c)
	for {
		msg, err := b.ReadBytes('\n')
		if err != nil {
			break
		}

		handleIncomin(msg, c)

		//fmt.Printf("%v:%s", c.RemoteAddr(), string(msg))
		//c.Write(msg)
		handleEvent(Event{Name: "CLIENT_MSG", Client: c, Msg: msg})
	}
}

func setupListenerTLS(ip string, sport string) {
	cert, err := tls.LoadX509KeyPair("secure/certs/server.pem", "secure/certs/server.key")
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	service := fmt.Sprintf("%s:%s", ip, sport)
	listener, err := tls.Listen("tcp", service, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
	}
	log.Print("server: listening")
	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Printf("server: accept: %s", err)
			fmt.Printf("%v\n", err)
			break
		}
		defer conn.Close()
		log.Printf("server: accepted from %s", conn.RemoteAddr())
		tlscon, ok := conn.(*tls.Conn)
		if ok {
			log.Print("ok=true")
			state := tlscon.ConnectionState()
			for _, v := range state.PeerCertificates {
				log.Print(x509.MarshalPKIXPublicKey(v.PublicKey))
			}
		}
		//go handleClient(conn)
		go handleConn(conn)
	}
}

func handleClient(conn net.Conn) {
	defer conn.Close()
	buf := make([]byte, 512)
	for {
		log.Print("server: conn: waiting")
		n, err := conn.Read(buf)
		if err != nil {
			if err != nil {
				log.Printf("server: conn: read: %s", err)
			}
			break
		}

		l := lex.PrvLexer{}
		l.Lex(buf)

		log.Printf("server: conn: echo %q\n", string(buf[:n]))
		n, err = conn.Write(buf[:n])
		log.Printf("server: conn: wrote %d bytes", n)

		if err != nil {
			log.Printf("server: write: %s", err)
			break
		}
	}
	log.Println("server: conn: closed")
}
