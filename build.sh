#!/bin/sh
# Simple script to build binaries in the current directory
go build ./cmd/app_sink
go build ./cmd/hub
go build ./cmd/app_rise
go build ./cmd/stompy
go build ./cmd/gobacker
# go build ./cmd/hub_sink # will fail for now, until we fix some stuff
