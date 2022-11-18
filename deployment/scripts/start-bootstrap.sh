#!/bin/bash

cd lotus || exit

tmux kill-session -t lotus
tmux new-session -d -s lotus
tmux send-keys "./lotus daemon --genesis=spacenet-genesis.car --profile=bootstrapper 2>&1" C-m
./lotus wait-api
./lotus net listen | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr
