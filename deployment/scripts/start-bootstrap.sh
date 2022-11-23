#!/bin/bash

cd lotus || exit

# Kill a potentially running instance of Lotus
tmux kill-session -t lotus
tmux new-session -d -s lotus

# Start the Lotus daemon and import the bootstrap key.
# Keeping the version with a custom genesis commented, in case we need to come back to it.
#tmux send-keys "./lotus daemon --genesis=spacenet-genesis.car --profile=bootstrapper --bootstrap=false 2>&1" C-m
tmux send-keys "./lotus daemon --profile=bootstrapper --bootstrap=false 2>&1" C-m
mkdir -p ~/.lotus/keystore && chmod 0700 ~/.lotus/keystore
./lotus-shed keyinfo import spacenet-libp2p-bootstrap1.keyinfo
echo '[Libp2p]
ListenAddresses = ["/ip4/0.0.0.0/tcp/1347"]' > ~/.lotus/config.toml
./lotus wait-api
./lotus net listen | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr

# Start the Faucet.
./lotus wallet import --as-default --format=json-lotus spacenet_faucet.key
cd ~/spacenet/faucet/ || exit
go build -o spacenet-faucet ./cmd || exit
tmux new-session -d -s faucet
tmux send-keys "export LOTUS_PATH=~/.lotus && ./spacenet-faucet 2>&1" C-m
