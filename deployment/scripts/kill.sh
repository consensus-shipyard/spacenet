#!/bin/bash

killall eudico
killall spacenet-faucet
tmux kill-server

sleep 3

killall -9 eudico
killall -9 spacenet-faucet
tmux kill-server

# Some of the above commands inevitably fail, as not all processes are running on all machines.
# This will prevent Ansible (through which this script is expected to be run) from complaining about it.
true