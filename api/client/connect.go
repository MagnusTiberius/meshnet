package client

import (
	"crypto/tls"
	"log"
	"net"
)

//ConfigTLS todo ...
type ConfigTLS struct {
	Addr      string
	ClientPEM string //"secure/certs/client.pem"
	ClientKey string //"secure/certs/client.key"
}

//NewTLS todo...
func NewTLS(cfg *ConfigTLS) net.Conn {
	cert, err := tls.LoadX509KeyPair(cfg.ClientPEM, cfg.ClientKey)
	if err != nil {
		log.Fatalf("server: loadkeys: %s", err)
		return nil
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}, InsecureSkipVerify: true}
	conn, err := tls.Dial("tcp", cfg.Addr, &config)
	if err != nil {
		log.Fatalf("client: dial: %s", err)
		return nil
	}
	log.Println("client: connected to: ", conn.RemoteAddr())
	/*
		state := conn.ConnectionState()
		for _, v := range state.PeerCertificates {
			fmt.Printf("PublicKey:\n")
			fmt.Println(x509.MarshalPKIXPublicKey(v.PublicKey))
			fmt.Printf("Subject:\n")
			fmt.Println(v.Subject)
		}
		log.Println("client: handshake: ", state.HandshakeComplete)
		log.Println("client: mutual: ", state.NegotiatedProtocolIsMutual)
	*/
	return conn
}
