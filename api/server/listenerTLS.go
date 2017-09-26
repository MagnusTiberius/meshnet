package server

import (
	"crypto/rand"
	"crypto/tls"
	"log"
	"net"
)

//Config todo ...
type Config struct {
	Addr string
	PEM  string //"secure/certs/client.pem"
	Key  string //"secure/certs/client.key"
}

//ListenerTLS todo ...
func ListenerTLS(cfg *Config) net.Listener {
	cert, err := tls.LoadX509KeyPair(cfg.PEM, cfg.Key)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
		return nil
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	config.Rand = rand.Reader
	listener, err := tls.Listen("tcp", cfg.Addr, &config)
	if err != nil {
		log.Fatalf("server: listen: %s", err)
		return nil
	}
	return listener
}
