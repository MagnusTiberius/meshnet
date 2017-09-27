package repo

import (
	"fmt"
	"net"

	"github.com/gomqtt/packet"
)

var (
	topicList map[string]*TopicItem
)

//TopicItem todo ...
type TopicItem struct {
	ConnList   map[string]net.Conn
	Messages   []*packet.Message
	CurrentMsg *packet.Message
}

//Bundle todo ...
type Bundle struct {
	TopicList map[string]*TopicItem
}

//NewBundle todo ...
func NewBundle() *Bundle {
	return &Bundle{
		TopicList: map[string]*TopicItem{},
	}
}

//Publish todo ...
func (t *Bundle) Publish(msg *packet.Message, conn net.Conn) {
	ti := t.TopicList[msg.Topic]
	if ti == nil {
		ti = &TopicItem{
			Messages:   []*packet.Message{msg},
			CurrentMsg: msg,
		}
		t.TopicList[msg.Topic] = ti
		return
	}
	ti.Messages = append(ti.Messages, msg)
	ti.CurrentMsg = msg
}

//Subscribe todo ...
func (t *Bundle) Subscribe(sub *packet.Subscription, conn net.Conn) {
	ti := t.TopicList[sub.Topic]
	if ti == nil {
		ti = &TopicItem{
			ConnList: map[string]net.Conn{fmt.Sprintf("%v", conn.RemoteAddr()): conn},
		}
		t.TopicList[sub.Topic] = ti
		return
	}
	ti.ConnList[fmt.Sprintf("%v", conn.RemoteAddr())] = conn
}
