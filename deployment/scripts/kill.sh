#!/bin/bash

killall lotus
killall mir-validator
killall spacenet-faucet
tmux kill-server

sleep 3

killall -9 lotus
killall -9 mir-validator
killall -9 spacenet-faucet
tmux kill-server
