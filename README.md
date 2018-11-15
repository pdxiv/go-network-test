# gonetworktest

## Introduction

Some tests of broadcast and UDP network functionality in Go.

## Prerequisites

### Program dependencies

- Requires Go version 1.11 or above to run/build.
- Currently only tested to work in Linux. (Possibly, the SO_REUSEPORT functionality won't work the same under Windows.)
- The autoconfig.pl relies on the Perl JSON module (available in Debian etc as `libjson-perl`).

### Building and running

To download and build:

```bash
go get github.com/pdxiv/gonetworktest
git clone https://github.com/pdxiv/gonetworktest
make
```

Network configuration settings are required before running. Settings are located in a `conf.json` file. Either edit this manually to adapt to your network settings, or run

```bash
./autoconfig.pl
```

### Performance

To make sure that we get enough performance in Linux, it's important that we remember to increase the default OS send and receive buffer size for all types of connections. Increasing it to something like 32 mb seems to work well for what we're trying to do here. It may be a good idea to increase the number of simultaneous open file handles to handle high load scenarios better.

```bash
sysctl -w net.core.rmem_max=33554432
sysctl -w net.core.wmem_max=33554432
sysctl -w net.core.rmem_default=33554432
sysctl -w net.core.wmem_default=33554432
sysctl -w fs.file-max=16777216
```

## Concepts

### Communication terminology

- Sink: Incoming communication to a service (think: "sink to").
- Rise: Outgoing communication from a service (think: "rise from").
- Hub: Central service handling all messages. Typically, your system will only have one active Hub.
- App: All other services. All App services communicate with each other over the Hub.
- Gob: Service for going back to historical sent data. Used during startup of a service and if messages are lost.

```text
           +-------+
  App Sink |       | App Rise
+---------->  App  +----------+
|          |       |          |
|          +-------+          |
|                             |
|          +-------+          |
| Hub Rise |       | Hub Sink |
+----------+  Hub  <----------+
           |       |
           +-------+
```

### App message handling

Since UDP doesn't guarantee message delivery, or message order, Apps receiving data from the hub need to have a mechanism for handling this. If one or more messages are lost, there is a gap in the sequence number, and the App will request the data with the missing sequence numbers from the "Gob" service. If a message with the same Hub sequence number has already been received, the message will be ignored.

```text
                     +------------+
                     |Get new     |
      +--------------+packets from<-------------+
      |              |Hub         |             |
      |              +-----^------+             |
      |                    |                    |
      |                    |Yes                 |
      |                    |                    |
+-----v------+       +-----+------+       +-----+------+
|Gap?        |   No  |Number too  |   No  |Send packet |
|            +------->low?        +------->to callee   |
|            |       |            |       |            |
+-----+------+       +------------+       +-----^------+
      |                                         |
      |Yes                                      |
      |                                         |
+-----v------+                            +-----+------+
|Get missing |                            |Send missing|
|packets from+---------------------------->packets to  |
|Gob         |                            |callee      |
+------------+                            +------------+
```

#### Communication protocols

The data fields all use network byte order (big-endian), when transmitted across the network.

AppRiseData carries information from an App to the Hub.

```golang
type AppRiseData struct {
    Type              uint16
    PayloadSize       uint16
    ID                uint64
    AppSequenceNumber uint64
    Payload           []byte
}
```

HubRiseData carries data from the Hub to Apps. Typically the payload contains one or more encapsulated App messages.

```golang
type HubRiseData struct {
    SessionID           uint64
    HubSequenceNumber   uint64
    NumberOfAppPayloads uint16 // If we put together several App msgs in one Hub
    Payload             []byte
}
```
