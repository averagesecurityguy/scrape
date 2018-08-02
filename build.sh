#!/bin/bash


INSTALL_DIR=/opt/scrape
GO_BIN=go

# Build our binary
GOPATH=/tmp/go "$GO_BIN" get -u github.com/asggo/store
GOPATH=/tmp/go "$GO_BIN" build -o bin/scrape
GOPATH=/tmp/go "$GO_BIN" build -o bin/kv github.com/asggo/store/src

# Stop Service
sudo service scrape stop

cp bin/scrape ${INSTALL_DIR}
cp bin/kv ${INSTALL_DIR}
# Start Service
sudo service scrape start
