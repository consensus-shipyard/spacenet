#!/bin/bash

cd lotus || exit

# Create a new wallet to be used by the validator
./lotus wallet new

# Initialize a new configuration for the mir validator.
# This will create mir-related config files in the $LOTUS_PATH directory.
./mir-validator config init

# Get the libp2p address of the local lotus node
lotus_listen_addr=$(./mir-validator config validator-addr | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/')

echo "${lotus_listen_addr}" > ~/.lotus/mir-validator-identity
