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
./build.sh
```

Network configuration settings are required before running. Settings are located in a `conf.json` file. Either edit this manually to adapt to your network settings, or run

```bash
./autoconfig.pl
```

### Performance

To make sure that we get enough performance in Linux, it's important that we remember to increase the default OS send and receive buffer size for all types of connections. Increasing it to something like 32 mb seems to work well for what we're trying to do here:

```bash
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

Since UDP doesn't guarantee message delivery, or message order, Apps receiving data from the hub need to have a mechanism for handling this. If one or more messages are lost, there is a gap in the sequence number, and the App will request the data with the missing sequence numbers from the "Gobacker" service. If a message with the same Hub sequence number has already been received, the message will be ignored.

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
|Gobacker    |                            |callee      |
+------------+                            +------------+
```

### Gobacker

#### Usage

The Gobacker service keeps track of what messages have been sent from the Hub. This is used in two situations:

- When an App is starting up and needs to read up on what messages have been sent to build internal state for a session. Typically, it broadcast "who has sequence number 0, for the latest SessionID (0xffffffffffffffff)", and when it gets a response from a Gobacker, it will connect to it via TCP and request messages with sequence numbers 0 to the largest possible sequence number (0xffffffffffffffff). The Gobacker will send as many packets as it has, and then closes down the connection, leaving the App to resume normal online operation. The App should keep track of what message sequence numbers have been sent out already for its' own AppID, so that it doesn't re-send messages to the Hub uselessly.
- When an App experiences a gap in sequence numbers from the Hub. The App then asks the Gobacker for the messages with the missing sequence numbers, for the current SessionID.
