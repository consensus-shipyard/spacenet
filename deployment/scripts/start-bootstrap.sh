#!/bin/bash

# Obtain bootstrap key file.
bootstrap_key="$1"
[ -n "${bootstrap_key}" ] || exit
shift

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
bootstrap_log_dir=~/spacenet-logs/bootstrap-$(date +%Y-%m-%d-%H-%M-%S_%Z)
faucet_log_dir=~/spacenet-logs/faucet-$(date +%Y-%m-%d-%H-%M-%S_%Z)
mkdir -p "$bootstrap_log_dir"
mkdir -p "$faucet_log_dir"

# Kill a potentially running instance of Lotus
tmux kill-session -t lotus
tmux new-session -d -s lotus

# Start the Lotus daemon and import the bootstrap key.
mkdir -p ~/.lotus/keystore && chmod 0700 ~/.lotus/keystore
./lotus-shed keyinfo import "${bootstrap_key}"
echo '[Libp2p]
ListenAddresses = ["/ip4/0.0.0.0/tcp/1347"]
[Chainstore]
  EnableSplitstore = true
[Chainstore.Splitstore]
  ColdStoreType = "discard"
' > ~/.lotus/config.toml
tmux send-keys "./eudico mir daemon --profile=bootstrapper --bootstrap=false 2>&1 | ./rotate-logs.sh ${bootstrap_log_dir} ${log_file_lines} ${max_archive_size}" C-m
./eudico wait-api
./eudico net listen | grep -vE '(/ip6/)|(127.0.0.1)|(/tcp/1347)' | grep -E '/ip4/.*/tcp/' > ~/.lotus/lotus-addr

# Start the Faucet.
./eudico wallet import --as-default --format=json-lotus spacenet_faucet.key
cd ~/spacenet/faucet/ || exit
go build -o spacenet-faucet ./cmd/faucet || exit
tmux new-session -d -s faucet
tmux send-keys "export LOTUS_PATH=~/.lotus && ./spacenet-faucet --web-host \"0.0.0.0:8000\" --web-allowed-origins \"*\" --web-backend-host \"https://spacenet.consensus.ninja/fund\" --filecoin-address=t1jlm55oqkdalh2l3akqfsaqmpjxgjd36pob34dqy --lotus-api-host=127.0.0.1:1234 2>&1 | ~/lotus/rotate-logs.sh ${faucet_log_dir} ${log_file_lines} ${max_archive_size}" C-m
