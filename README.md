# go-network-test

## Introduction
Some tests of broadcast and UDP network functionality in Go.

## Prerequisites

### Program dependencies
- Requires Go version 1.10 or above to run/build.
- Relies on `github.com/libp2p/go-reuseport` that sets `SO_REUSEPORT` and `SO_REUSEADDR`, which is functionality that won't be available out of the box before Go version 1.11.
- Currently only works in Linux, because of the dependency on `github.com/libp2p/go-reuseport`.
- The autoconfig.pl relies on the Perl JSON module (available in Debian etc as `libjson-perl`).
### Building and running
To download and build:
```
git clone https://github.com/pdxiv/gonetworktest
go get github.com/pdxiv/gonetworktest
./build.sh
```
Network configuration settings are required before running. Settings are located in a `conf.json` file. Either edit this manually to adapt to your network settings, or run
```
./autoconfig.pl
```
### Performance
To make sure that we get enough performance in Linux, it's important that we remember to increase the default OS send and receive buffer size for all types of connections. Increasing it to something like 32 mb seems to work well for what we're trying to do here:
```
sysctl -w net.core.rmem_max=33554432
sysctl -w net.core.wmem_max=33554432
sysctl -w net.core.rmem_default=33554432
sysctl -w net.core.wmem_default=33554432
```
## Concepts
### Communication terminology
- Sink: Incoming communication to a service (think, "sink to").
- Rise: Outgoing communication from a service (think, "rise from").
- Hub: Central service handling all messages. Typically, your system will only have one active Hub.
- App: All other services. All App services communicate with each other over the Hub.
```
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
Since UDP doesn't guarantee message delivery, or message order, Apps receiving data from the hub need to have a mechanism for handling this. If one or more messages are lost, there is a gap in the sequence number, and the App will request the data with the missing sequence numbers from the "Gobacker" service. If a message with the same Hub sequence number has already been received, the message will be ignored.
```
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
|Gobacker    |                            |callee      |
+------------+                            +------------+
```
