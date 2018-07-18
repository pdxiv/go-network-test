#!/usr/bin/perl
# Very simple "happy path" script to attempt to create a config file automatically
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

print $file_handle encode_json($json_data);

close $file_handle;

sub default_route_interface {
    my @route_info = `route`;
    foreach (@route_info) {
        if (/^default.*?(\w+)$/) {
            return $1;
        }
    }
    return q{};
}

sub broadcast_for_interface {
    my $interface         = shift;
    my @interface_info    = `ifconfig`;
    my $current_interface = q{};
    my %interface_broadcast;
    foreach my $interface_line (@interface_info) {
        if ( $interface_line =~ /^([^: \n]+)/ ) {
            $current_interface = $1;
        }
        if ( $interface_line =~ /(?: Bcast:| broadcast )\s*(\S+)/ ) {
            $interface_broadcast{$current_interface} = $1;
        }
    }
    return $interface_broadcast{$interface};
}
