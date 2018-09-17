#!/bin/bash


INSTALL_DIR=/opt/scrape
GO_BIN=go
USER=scrape

# Build our binary
GOPATH=/tmp/go "$GO_BIN" get -u github.com/boltdb/bolt
GOPATH=/tmp/go "$GO_BIN" get -u github.com/gorilla/mux
GOPATH=/tmp/go "$GO_BIN" build -o bin/scrape
GOPATH=/tmp/go "$GO_BIN" build -o bin/view view/view.go view/store.go

# Stop Service
sudo service scrape stop

sudo cp bin/scrape ${INSTALL_DIR}
sudo cp bin/view ${INSTALL_DIR}
sudo cp -R web ${INSTALL_DIR}
sudo chown -R $USER:$USER ${INSTALL_DIR}

# Start Service
sudo service scrape start
