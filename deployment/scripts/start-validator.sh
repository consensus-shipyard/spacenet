#!/bin/bash

cd lotus || exit

tmux kill-session -t mir-validator
tmux new-session -d -s mir-validator
tmux send-keys "./mir-validator run --nosync 2>&1" C-m
