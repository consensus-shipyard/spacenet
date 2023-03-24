#!/bin/bash

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
log_dir=~/spacenet-logs/validator-$(date +%Y-%m-%d-%H-%M-%S_%Z)
mkdir -p "$log_dir"

# Kill a potentially running validator.
tmux kill-session -t mir-validator
tmux new-session -d -s mir-validator

# Start the Mir validator.
tmux send-keys "LOTUS_PATH=/home/ubuntu/.lotus ./eudico mir validator run --nosync --max-block-delay=15s 2>&1 | ./rotate-logs.sh ${log_dir} ${log_file_lines} ${max_archive_size}" C-m
