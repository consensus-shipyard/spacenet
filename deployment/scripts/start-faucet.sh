#!/bin/bash

# Obtain number of lines per log file.
log_file_lines="$1"
[ "${log_file_lines}" -gt 0 ] || exit
shift

# Obtain maximal log archive size (in bytes).
max_archive_size=$1
[ "${max_archive_size}" -gt 0 ] || exit
shift

# Create log directory.
faucet_log_dir=~/spacenet-logs/faucet-$(date +%Y-%m-%d-%H-%M-%S_%Z)
mkdir -p "$faucet_log_dir"

# Kill a potentially running instance of the faucet.
tmux kill-session -t faucet

# Import faucet key to Eudico (he faucet's address has the coins that will be distributed).
cd lotus || exit
./eudico wallet import --as-default --format=json-lotus spacenet_faucet.key

tag=$(git describe --tags 2>/dev/null || echo "unk-$(git rev-parse --short=10 HEAD)")
flags="-X=github.com/filecoin-project/faucet/pkg/version.gittag=${tag}"

# Start the Faucet.
cd ~/spacenet/faucet/ || exit
go build -o spacenet-faucet -ldflags "$flags" ./cmd/faucet || exit
tmux new-session -d -s faucet
tmux send-keys "export LOTUS_PATH=~/.lotus && ./spacenet-faucet --web-host \"0.0.0.0:8000\" --web-allowed-origins \"*\" --web-backend-host \"https://spacenet.consensus.ninja/fund\" --filecoin-address=t1jlm55oqkdalh2l3akqfsaqmpjxgjd36pob34dqy --lotus-api-host=127.0.0.1:1234 2>&1 | ~/lotus/rotate-logs.sh ${faucet_log_dir} ${log_file_lines} ${max_archive_size}" C-m
