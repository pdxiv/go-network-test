#!/usr/bin/perl
use strict;
use warnings;
our $VERSION = '1.0.0';
use Socket;

# create a socket
socket TO_SERVER, PF_INET, SOCK_STREAM, getprotobyname 'tcp';

my $remote_host = '127.0.0.1';
my $remote_port = '9996';

# build the address of the remote machine
my $internet_addr = inet_aton($remote_host)
  or die "Couldn't convert $remote_host into an Internet address: $!\n";
my $paddr = sockaddr_in( $remote_port, $internet_addr );

# connect
connect TO_SERVER, $paddr or die "Couldn't connect to $remote_host:$remote_port : $!\n";

# ... do something with the socket
# print TO_SERVER "Why don't you call me anymore?\n";
print TO_SERVER pack "H*", 0x1337;
# TO_SERVER->flush();

# and terminate the connection when we're done
close TO_SERVER;
