# go-network-test
Some tests of broadcast and udp network functionality in Go. Currently only works in Linux, because of a dependency on a hack that sets SO_REUSEPORT and SO_REUSEADDR, which is functionality that wont' be available in Go until version 1.11.
