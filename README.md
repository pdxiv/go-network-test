# go-network-test
Some tests of broadcast and UDP network functionality in Go. Currently only works in Linux, because of a dependency on a hack that sets `SO_REUSEPORT` and `SO_REUSEADDR`, which is functionality that won't be available out of the box before Go version 1.11.

To make sure that we get enough performance in Linux, it's important that we remember to increase the default OS send and receive buffer size for all types of connections. Increasing it to something like 32 mb seems to work well for what we're trying to do here:
```
sysctl -w net.core.rmem_max=33554432
sysctl -w net.core.wmem_max=33554432
sysctl -w net.core.rmem_default=33554432
sysctl -w net.core.wmem_default=33554432
```
