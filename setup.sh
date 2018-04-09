#!/bin/bash


INSTALL_DIR=/opt/scrape
USER=scrape

# Install Dependencies
sudo apt install golang

# Build our binary
go build

# Install
mkdir ${INSTALL_DIR}
mkdir ${INSTALL_DIR}/data
mkdir ${INSTALL_DIR}/var
mkdir ${INSTALL_DIR}/log

cp scrape ${INSTALL_DIR}
chown -R ${USER}:${USER} ${INSTALL_DIR}

# Configure Service
sudo cp service.sh /etc/init.d/scrape
sudo chmod +x /etc/init.d/scrape
sudo update-rc.d scrape defaults

# Start Service
sudo service scrape start
