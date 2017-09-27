package command

import (
	"github.com/MagnusTiberius/packet"
)

//NewMessage returns a new packet message
func NewMessage() *packet.Message {
	return &packet.Message{}
}
