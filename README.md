# meshnet
Secure Mesh Topology

Decentralized distributed MQTT broker and client running on TCP and TLS (key and pem). Use private and public keys to implement trust connection in an untrusted public network. Data is encrypted using TLS and joining a network require public keys.


Meshnet is an implementation of MQTT running on TCP. It supports TLS using SSL certificates.

Protocol:
```
  MQTT
   http://mqtt.org/documentation
```

Status:
```
   9/25/2017 - Implemented connect and connect-ack.
   9/27/2017 - Implemented publish, subscribe.
               Implemented message dispatcher.
               Implemented removal of closed connections.
   9/29/2017 - Implemented pingreq, pingresp, disconnect
```

Dependencies:

1. go get github.com/MagnusTiberius/packet
