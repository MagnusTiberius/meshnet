# meshnet
Secure Mesh Topology

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
```

Dependencies:

1. go get github.com/MagnusTiberius/packet
