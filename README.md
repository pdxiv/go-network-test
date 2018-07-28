# go-network-test

## Introduction
Some tests of broadcast and UDP network functionality in Go.

## Prerequisites

### Program dependencies
- Requires Go version 1.10 or above to run/build.
- Relies on `github.com/libp2p/go-reuseport` that sets `SO_REUSEPORT` and `SO_REUSEADDR`, which is functionality that won't be available out of the box before Go version 1.11.
- Currently only works in Linux, because of the dependency on `github.com/libp2p/go-reuseport`.
- The autoconfig.pl relies on the Perl JSON module (available in Debian etc as `libjson-perl`).

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
- Seq: Central service handling all messages. Typically, your system will only have one active Seq.
- App: All other services. All App services communicate with each other over the Seq.
```
           +-------+
  App Sink |       | App Rise
+---------->  App  +----------+
|          |       |          |
|          +-------+          |
|                             |
|          +-------+          |
| Seq Rise |       | Seq Sink |
+----------+  Seq  <----------+
           |       |
           +-------+
```
