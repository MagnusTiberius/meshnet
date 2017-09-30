# meshnet
Secure Mesh Topology

Decentralized distributed MQTT broker and client running on TCP and TLS (key and pem). Use private and public keys to implement trust connection in an untrusted public network. Data is encrypted using TLS and joining a network require public keys.


Meshnet is an implementation of MQTT running on TCP. It supports TLS using SSL certificates.

##Protocol
```
  MQTT
   http://mqtt.org/documentation
```

##Status
```
   9/25/2017 - Implemented connect and connect-ack.
   9/27/2017 - Implemented publish, subscribe.
               Implemented message dispatcher.
               Implemented removal of closed connections.
   9/29/2017 - Implemented pingreq, pingresp, disconnect
               Implemented packet callbacks
               Setup two sensor publishers and one mobile device subscriber
   9/30/2017 - Implemented unsubscribe and unsubscribe-ack
```

##Dependencies

1. go get github.com/MagnusTiberius/packet

##Example | Demo

```
   Broker: 
       - See sample broker program
       - usage: go run ./server/main.go
```

```
   Publishers:
       - Two demo programs are provided for pushing sensor data to the broker.
       - usage: go run ./agentsensor1/main.go
       - usage: go run ./agentsensor2/main.go
```

```
   Subscriber:
       - One demo program for a subscriber.
       - usage: go run ./agentmobiledevice1/main.go
```

