#!/bin/bash

# Obtain bootstrap node's address as a parameter.
bootstrap_addr="$1"
[ -n "$bootstrap_addr" ] || exit
shift

# Obtain number of lines per log file.
log_file_lines="$1"
[ "${log_file_lines}" -gt 0 ] || exit
shift

# Make sure that the maximal log archive size (in bytes) has been properly specified.
max_archive_size=$1
[ "${max_archive_size}" -gt 0 ] || exit
shift


cd lotus || exit

# Create log directory
log_dir=~/spacenet-logs/daemon-$(date +%Y-%m-%d-%H-%M-%S_%Z)
mkdir -p "$log_dir"

# Kill a potentially running instance of Lotus
tmux kill-session -t lotus
tmux new-session -d -s lotus

# Start the Lotus daemon and import the bootstrap key.
# Keeping the version with a custom genesis commented, in case we need to come back to it.
#tmux send-keys "./lotus daemon --genesis=spacenet-genesis.car --bootstrap=true --mir-validator 2>&1" C-m
tmux send-keys "./lotus daemon --bootstrap=true --mir-validator 2>&1 | ./rotate-logs.sh ${log_dir} ${log_file_lines} ${max_archive_size}" C-m
./lotus wait-api
./lotus net connect "$bootstrap_addr"
./lotus net listen | grep -vE '(/ip6/)|(127.0.0.1)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr
