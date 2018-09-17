#!/bin/bash


INSTALL_DIR=/opt/scrape
GO_BIN=go

# Build our binary
GOPATH=/tmp/go "$GO_BIN" get -u github.com/boltdb/bolt
GOPATH=/tmp/go "$GO_BIN" get -u github.com/gorilla/mux
GOPATH=/tmp/go "$GO_BIN" build -o bin/scrape
GOPATH=/tmp/go "$GO_BIN" build -o bin/view view/view.go view/store.go

# Stop Service
sudo service scrape stop

cp bin/scrape ${INSTALL_DIR}
cp bin/view ${INSTALL_DIR}
cp -R web ${INSTALL_DIR}

# Start Service
sudo service scrape start
