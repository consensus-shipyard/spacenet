#!/bin/bash

# Obtain bootstrap node's address as a parameter
bootstrap_addr="$1"
if [ -z "$bootstrap_addr" ]; then
  exit 1
fi

cd lotus || exit

# Kill a potentially running instance of Lotus
tmux kill-session -t lotus
tmux new-session -d -s lotus

# Start the Lotus daemon and import the bootstrap key.
# Keeping the version with a custom genesis commented, in case we need to come back to it.
#tmux send-keys "./lotus daemon --genesis=spacenet-genesis.car --bootstrap=true --mir-validator 2>&1" C-m
tmux send-keys "./lotus daemon --bootstrap=true --mir-validator 2>&1" C-m
./lotus wait-api
./lotus net connect "$bootstrap_addr"
./lotus net listen | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr
