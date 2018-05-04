#!/bin/bash


INSTALL_DIR=/opt/scrape
GO_BIN=go

# Build our binary
GOPATH=/tmp/go "$GO_BIN" get github.com/asggo/store
GOPATH=/tmp/go "$GO_BIN" build -o bin/scrape
GOPATH=/tmp/go "$GO_BIN" build -o bin/search ./search/

# Stop Service
sudo service scrape stop

cp bin/scrape ${INSTALL_DIR}
cp bin/search ${INSTALL_DIR}

# Start Service
sudo service scrape start
