package main

import "fmt"

func handleEvent(e Event) {
	//fmt.Println(e.Name)
	switch e.Name {
	case "CONNECT_EVENT":
		for _, c := range connPool {
			conn := c
			addr := fmt.Sprintf("%v", e.Client.RemoteAddr())
			if fmt.Sprintf("%v", conn.RemoteAddr()) != addr {
			}
			msgConnect := []byte(fmt.Sprintf("%s has connected\n", addr))
			conn.Write(msgConnect)
		}
	case "CLIENT_MSG":
		for _, c := range connPool {
			addr := fmt.Sprintf("%v", e.Client.RemoteAddr())
			conn := c
			if fmt.Sprintf("%v", conn.RemoteAddr()) != addr {
				msg := fmt.Sprintf("ECHO>>%v:%s", e.Client.RemoteAddr(), string(e.Msg))
				conn.Write([]byte(msg))
			}
		}
	case "RELAY_MSG":
		for _, c := range connPool {
			client := c
			msg := fmt.Sprintf("RELAY>>%v:%s\n", e.Client.RemoteAddr(), string(e.Msg))
			client.Write([]byte(msg))
		}
	default:
	}
	//relayEvent(e)
}
