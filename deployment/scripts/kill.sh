#!/bin/bash

killall eudico
killall mir-validator
killall spacenet-faucet
tmux kill-server

sleep 3

killall -9 eudico
killall -9 mir-validator
killall -9 spacenet-faucet
tmux kill-server
