#!/bin/bash

killall -9 lotus
killall -9 mir-validator
killall -9 spacenet-faucet
tmux kill-server
