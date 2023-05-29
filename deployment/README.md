# Deploying Spacenet

We use [Ansible](https://www.ansible.com/) to deploy Spacenet nodes (and whole networks).
The set of machines on which to deploy Spacenet must be defined in an Ansible inventory file
(we use `hosts` as an example inventory file name in this document, but any other file name is also allowed).
The inventory file must contain 2 host groups: `bootstrap` (with 1 host) and `validators` (with all validator hosts).
An example host file looks as follows.
```
[bootstrap]
198.51.100.0

[validators]
198.51.100.1
198.51.100.2
198.51.100.3
198.51.100.4
```

The Spacenet deployment can be managed using the provided Ansible playbooks.
To run a playbook, install Ansible and execute the following command
```shell
ansible-playbook -i hosts <playbook.yaml> ...
```
with `hosts` being the Ansible inventory and `<playbook.yaml>` one of the provided playbooks.
Additional playbooks can be specified in the same command and will be executed in the given sequence.
A reference of the provided deployment playbooks is given at the end of this document.

## Choosing deployment targets

Running the command above applies the playbooks to their default targets,
assuming all nodes in the inventory are part of Spacenet.
To target specific nodes from the inventory, the `nodes` Ansible variable can be used
through specifying an additional parameter `--extra-vars "nodes='<alternative_targets>'"`.
For example the following commands, respectively,
only set up the bootstrap node and only kill the validators `198.51.100.3` and `198.51.100.4`.
```shell
ansible-playbook -i hosts setup.yaml --extra-vars "nodes=bootstrap"
ansible-playbook -i hosts kill.yaml --extra-vars "nodes='198.51.100.3 198.51.100.4'"
```

## Ansible parallelism

By default, ansible communicates with 5 remote nodes at a time.
This is fine for, say 4 validators and 1 bootstrap, but as soon as more nodes are involved,
it slows down the deployment significantly.
To increase the number of parallel ansible connections, use the `--forks` command-line argument.

```shell
ansible-playbook -i hosts --forks 10 <playbook.yaml> ...
```

## System requirements and configuration

- Ansible installed on the local machine.
- Python 3 (command `python3`) installed on the local machine.
- Ubuntu 22.04 on all remote machines (might easily work with other systems, but was tested with this one).
- Sudo access without password on remote machines.
- SSH access to remote machines without password

The file [group_vars/all.yaml](group_vars/all.yaml) contains some configuration parameters
(e.g. the location of the SSH key to use for accessing remote machines) documented therein.

### Potential issue on Ubuntu 22.04

While testing the deployment on Amazon EC2 virtual machines,
we noticed that installing dependencies on remote machines (performed by the `setup.yaml` playbook) sometimes failed.
The issue and solution has been
[described here](https://askubuntu.com/questions/1431786/grub-efi-amd64-signed-dependency-issue-in-ubuntu-22-04lts).
To apply the work-around, the `custom-script.yaml` playbook can be used.
If necessary, copy the following line
```shell
sudo apt --only-upgrade install grub-efi-amd64-signed
```
in the [scripts/custom.sh](scripts/custom.sh) file run
```shell
ansible-playbook -i hosts custom-script.yaml
```

## Deploying a fresh instance of Spacenet

To deploy an instance of Spacenet, first create an inventory file (called `hosts` in this example)
and populate it with IP addresses of machines that should run Spacenet as described above.

The following steps must be executed to deploy Spacenet:
1. Install necessary packages,
   clone the Spacenet client (Lotus) code and compile it on the remote machines (`setup.yaml`).
2. Start the bootstrap node (`start-bootstrap.yaml`)
3. Start the Lotus daemons on validator nodes (`start-daemons.yaml`)
4. Start the Mir validator processes on validator nodes (`start-validators.yaml`)

These steps are automated for convenience in the `deploy-new.yaml` playbook.
Thus, to deploy a fresh instance of Spacenet, simply run
```shell
ansible-playbook -i hosts deploy-new.yaml
```

## Rolling updates

When the Lotus code, the validator code, or the Mir code are updated,
the update can be rolled out to the running deployment, as long as the protocol remains the same.
For this, we provide the `rolling-update.sh` script.
This script performs a rolling update of selected nodes from an Ansible inventory and is invoked as follows.

```shell
./rolling-update hosts 198.51.100.1 [198.51.100.2 [...]]
```

The first argument must be an Ansible inventory file that contains all the node arguments that follow.
This script updates (fetches the code, recompiles it, and restarts the node, using the update-nodes.yaml)
the nodes one by one, always waiting for a node to catch up with the others
and only then proceeding to updating the next one.

## Provided deployment playbooks

### `clean.yaml`

Kills running Lotus daemon and Mir validator and deletes their associated state.
Does not touch the code and binaries.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `connect-daemons.yaml`

Connects all Lotus daemons to each other. This is required for the nodes to be able to sync their state.
It assumes the daemons and the bootstrap are up and running (but not necessarily the validators)

Applies to all hosts (including bootstrap).

### `custom-script.yaml`

Runs the scripts/custom.sh script. This is meant as a convenience tool for executing ad-hoc scripts.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `deep-clean.yaml`

Performs deep cleaning of the host machines.
Runs clean.yaml and, in addition, deletes the cloned repository with the lotus code and binaries.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `deploy-current.yaml`

Deploys the bootstrap and the validators using existing binaries.
This playbook still cleans the lotus daemon state, but neither updates nor recompiles the code.
Performs the state cleanup by running clean.yaml,
potentially producing and ignoring some errors, if nothing is running on the hosts - this is normal.

The nodes variable must not be set, as this playbook must distinguish between different kinds of nodes
(such as bootstrap and validators).

### `deploy-new.yaml`

Deploys the whole system from scratch.
Performs a deep clean by running deep-clean.yaml
(potentially producing and ignoring some errors, if nothing is running on the hosts - this is normal)
and sets up a new Spacenet deployment.

The nodes variable must not be set, as this playbook must distinguish between different kinds of nodes
(such as bootstrap and validators).

### `fetch-logs.yaml`

Fetches logs from all hosts and stores them in the `fetched-logs` directory (one sub-directory per host).

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `kill.yaml`

Kills running Lotus daemon and Mir validator.
Does not touch their persisted state or the code and binaries.
Reports but ignores errors, so it can be used even if the processes to be killed are not running.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `restart-validators.yaml`

Restarts a given set of validators.
For safety, does NOT default to restarting all validators
and the set of hosts to restart must be explicitly given using --extra-vars "nodes=..."

Note that this playbook always affects all hosts, regardless of the value of the nodes variable.
This is due to the necessity of reconnecting all daemons to the restarted one.

### `setup.yaml`

Sets up the environment for running the Lotus daemon and validator.
This includes installing the necessary packages, fetching the Lotus code, and compiling it.
It does not start any nodes. See start-* and deploy-*.yaml for starting the nodes.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-bootstrap.yaml`

Starts the bootstrap node and downloads its identity to localhost (using setup.yaml).

Applies to the bootstrap host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-faucet.yaml`

Starts the faucet service that can be used to distribute coins.
Assumes that the bootstrap node is up and running (see start-bootstrap.yaml).

Applies to the first bootstrap host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-daemons.yaml`

Starts the Lotus daemons and creates connections among them and to the bootstrap node.
Assumes that the bootstrap node is up and running (see start-bootstrap.yaml).

Applies to the validator host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-monitoring.yaml`

Starts the health monitoring service that can be used to check the status of the system.
Assumes that all the nodes (bootstraps, daemons, and validators) are up and running.

Applies to the validator host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-validators.yaml`

Starts the Mir validators.
Assumes that the Lotus daemons are up and running (see start-daemons.yaml).

Applies to the validator host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `update-nodes.yaml`

Updates a given set of validators by fetching the configured code, recompiling it, and restarting the validators.
After the update, waits until the nodes sync with the state of a bootstrap node and only then returns.

For safety, does NOT default to restarting all validators
and the set of hosts to restart must be explicitly given using --extra-vars "nodes=..."

Note that this playbook always affects all hosts, regardless of the value of the nodes variable.
This is due to the necessity of reconnecting all daemons to the restarted one.

### `status.yaml`

Gets the status of the Eudico daemons and the chain head.