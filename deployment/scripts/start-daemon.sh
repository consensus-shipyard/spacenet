#!/bin/bash

# Obtain bootstrap node's address as a parameter
bootstrap_addr="$1"
if [ -z "$bootstrap_addr" ]; then
  exit 1
fi

cd lotus || exit

tmux kill-session -t lotus
tmux new-session -d -s lotus
tmux send-keys "./lotus daemon --genesis=spacenet-genesis.car --bootstrap=false --mir-validator 2>&1" C-m
./lotus wait-api
./lotus net connect "$bootstrap_addr"
./lotus net listen | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr
