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
	LastSent   int64
}

//Bundle todo ...
type Bundle struct {
	TopicList map[string]*TopicItem
}

//NewBundle todo ...
func NewBundle() *Bundle {
	return &Bundle{
		TopicList: map[string]*TopicItem{"test*": &TopicItem{}},
	}
}

//Publish todo ...
func (t *Bundle) Publish(msg *packet.Message, conn net.Conn) {
	addr := fmt.Sprintf("%v", conn.RemoteAddr())
	fmt.Printf("Repo.Publish: %v \n", addr)
	ti, ok := t.TopicList[msg.Topic]
	if !ok {
		ti = &TopicItem{
			Messages:   []*packet.Message{msg},
			CurrentMsg: msg,
			LastSent:   -1,
		}
		t.TopicList[msg.Topic] = ti
		return
	}
	ti.Messages = append(ti.Messages, msg)
	ti.CurrentMsg = msg
}

//Subscribe todo ...
func (t *Bundle) Subscribe(sub *packet.Subscription, conn net.Conn) {
	addr := fmt.Sprintf("%v", conn.RemoteAddr())
	fmt.Printf("Repo.Subscribe: %v \n", addr)
	if t.TopicList == nil {
		t.TopicList = make(map[string]*TopicItem)
	}
	ti, ok := t.TopicList[sub.Topic]
	if !ok {
		//fmt.Printf("Repo.Subscribe 2: %v : %v \n", addr, sub.Topic)
		t.TopicList[sub.Topic] = &TopicItem{
			ConnList: map[string]net.Conn{addr: conn},
		}
		return
	}
	if ti.ConnList == nil {
		ti.ConnList = make(map[string]net.Conn)
	}
	_, ok = ti.ConnList[addr]
	if !ok {
		//fmt.Printf("Repo.Subscribe 3: %v : %v \n", addr, sub.Topic)
		ti.ConnList[addr] = conn
		return
	}
	//fmt.Printf("Repo.Subscribe 4: %v \n", addr)
	ti.ConnList[addr] = conn
}
