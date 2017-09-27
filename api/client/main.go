package client

import (
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
