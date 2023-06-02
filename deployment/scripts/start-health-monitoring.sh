#!/bin/bash

# Obtain number of lines per log file.
log_file_lines="$1"
[ "${log_file_lines}" -gt 0 ] || exit
shift

# Obtain maximal log archive size (in bytes)
max_archive_size=$1
[ "${max_archive_size}" -gt 0 ] || exit
shift

cd lotus || exit

# Create log directories
health_log_dir=~/spacenet-logs/health-$(date +%Y-%m-%d-%H-%M-%S_%Z)
mkdir -p "$health_log_dir"

gittag=$(git tag -l --sort=-creatordate | head -n 1 || echo "unk")
: ${gittag:="unk"}
githash=$(git rev-parse --short=8 HEAD)
flags="-X=github.com/filecoin-project/faucet/pkg/version.gittag=${gittag}"
flags+=" -X=github.com/filecoin-project/faucet/pkg/version.githash=${githash}"

# Start the Hello service.
cd ~/spacenet/faucet/ || exit
go build -o spacenet-health -ldflags "$flags" ./cmd/health || exit
tmux new-session -d -s health
tmux send-keys "export LOTUS_PATH=~/.lotus && ./spacenet-health --web-host \"0.0.0.0:9000\" --lotus-api-host=127.0.0.1:1234 2>&1 | ~/lotus/rotate-logs.sh ${health_log_dir} ${log_file_lines} ${max_archive_size}" C-m
