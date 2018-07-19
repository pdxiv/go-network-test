#!/usr/bin/perl
# Very simple "happy path" script to attempt to create a config file
# automatically.
# On an Ubuntu/Debian based system, this Perl script may need libjson-perl.
use strict;
use warnings;
our $VERSION = '0.1.0';
use utf8;
use JSON;

my $broadcast_address = broadcast_for_interface( default_route_interface() );

my $json_data = {
    'AppRiseAddress'       => "$broadcast_address:9998",
    'AppSinkAddress'       => '0.0.0.0:9999',
    'SequencerRiseAddress' => "$broadcast_address:9999",
    'SequencerSinkAddress' => '0.0.0.0:9998',
};
open my $file_handle, q{>}, 'conf.json';

print {$file_handle} encode_json($json_data);

close $file_handle;

sub default_route_interface {
    my $ip_route_default_interface = qr/^default.* dev (\S+)/;
    foreach (`ip route`) {
        if (/${ip_route_default_interface}/) {
            return $1;
        }
    }
    return q{};
}

sub broadcast_for_interface {
    my $interface_to_find = shift;
    my $ip_addr_interface = qr/^\d+:\s+([^:]+)/;
    my $ip_addr_broadcast = qr/ brd ([\d.]+)/;
    my $interface;
    my %broadcast_for_interface;
    foreach my $line (`ip addr`) {
        if ( $line =~ /${ip_addr_interface}/ ) {
            $interface = $1;
        }
        if ( $line =~ /${ip_addr_broadcast}/ ) {
            $broadcast_for_interface{$interface} = $1;
        }
    }
    if ( exists $broadcast_for_interface{$interface_to_find} ) {
        return $broadcast_for_interface{$interface_to_find};
    }
    return q{};
}
