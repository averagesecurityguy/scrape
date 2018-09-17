#!/bin/bash

INSTALL_DIR=/opt/scrape
USER=scrape
GO_BIN=go

# Create service user account
# Comment out these lines to prevent creating a new user account.
echo "Adding $USER account."
pass=$(head -c 12 /dev/urandom | base64)
useradd -s /bin/bash $USER
echo $USER:$pass | chpasswd
echo "User account $USER created with password $pass."

# Install Dependencies
sudo apt install golang

# Create Install directory
mkdir ${INSTALL_DIR}
mkdir ${INSTALL_DIR}/data
mkdir ${INSTALL_DIR}/var
mkdir ${INSTALL_DIR}/log
chown -R ${USER}:${USER} ${INSTALL_DIR}

# Configure Service
sudo cp service.sh /etc/init.d/scrape
sudo chmod +x /etc/init.d/scrape
sudo update-rc.d scrape defaults

# Compile and install binaries
./build.sh
