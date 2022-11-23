#!/bin/bash

# This script performs a rolling update of selected nodes from an Ansible inventory.
#
# usage: rolling-update inventory_file node1 [node2 [...]]
#
# The first argument must be an Ansible inventory file that contains all the node arguments that follow.
# This script updates (fetches the code, recompiles it, and restarts the node, using the update-nodes.yaml)
# the nodes one by one, always waiting for a node to catch up with the others
# and only then proceeding to updating the next one.

inventory="$1"
shift

while [ -n "$1" ]; do
  echo -e "\n========================================"
  echo "Updating node: $1"
  echo -e "========================================\n"
  ansible-playbook -i "$inventory" update-nodes.yaml --extra-vars "nodes=$1" || exit
  shift
done