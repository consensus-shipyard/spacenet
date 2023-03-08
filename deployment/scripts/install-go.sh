#!/bin/bash

# THIS IS A DIRTY HACK
# It is tested on Ubuntu Linux 22.04.
# No guarantees for other systems.
# Anyway, when the required and available versions of Go change (as they do all the time),
# this script will probably not be necessary and we'll be able to install Go using apt.
# Only when writing this, the apt version of Go was outdated.

wget https://go.dev/dl/go1.19.7.linux-amd64.tar.gz
sudo rm -rf /usr/local/go && sudo tar -C /usr/local -xzf go1.19.7.linux-amd64.tar.gz
sudo rm /usr/bin/go
sudo ln -s /usr/local/go/bin/go /usr/bin/go
sudo rm /usr/lib/go
sudo ln -s /usr/local/go /usr/lib/go
sudo rm /usr/share/go
sudo ln -s /usr/local/go /usr/share/go
