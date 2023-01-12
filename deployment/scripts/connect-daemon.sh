#!/bin/bash

cd lotus || exit

while IFS="" read -r addr || [ -n "$addr" ]; do         # Read all lotus addresses from the provided file, liney by line
  if [ "$addr" != "$(cat ../.lotus/lotus-addr)" ]; then # Skip own address
    ./eudico net connect "$addr" || exit                 # Connect to each other address
  fi
done < ../.lotus/lotus-addrs
