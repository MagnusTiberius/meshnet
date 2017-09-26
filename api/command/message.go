package command

import (
	"github.com/gomqtt/packet"
)

//NewMessage returns a new packet message
func NewMessage() *packet.Message {
	return &packet.Message{}
}
