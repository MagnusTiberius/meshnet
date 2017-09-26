package client

import (
	"bufio"
	"log"
	"net"

	"github.com/MagnusTiberius/meshnet/api/command"
	"github.com/gomqtt/packet"
)

//Config todo ...
type Config struct {
	Listener       net.Listener
	HandleIncoming func(reply []byte, conn net.Conn)
}

//NewClient todo ...
func NewClient() *Config {
	return &Config{}
}

//Start todo ...
func (c *Config) Start() {

}

//Connect todo ...
func (c *Config) Connect(packet *packet.ConnectPacket, conn net.Conn) {
	command.Connect(packet, conn)
}

//HandleReceive todo ..
func (c *Config) HandleReceive(conn net.Conn) {
	//reply := make([]byte, 4096)
	br := bufio.NewReader(conn)
	for {
		//n, err := conn.Read(reply)
		msg, err := br.ReadBytes('\n')
		if err != nil {
			log.Fatalf("client: write: %s", err)
		}
		//log.Printf("client: read %q (%d bytes)", string(reply[:n]), n)
		c.HandleIncoming(msg, conn)
	}
}
