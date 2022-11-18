# Deploying Spacenet

We use the [Ansible](https://www.ansible.com/) tool to deploy Spacenet nodes (and whole networks).
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
A reference of the provided deployment playbooks is provided at the end of this document.

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

## System requirements and configuration

- Ansible installed on the local machine.
- Ubuntu 22.04 on all remote machines (might easily work with other systems, but was tested with this one).
- Sudo access without passowrd on remote machines.
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
2. Generate the genesis block for the network and distribute it to all nodes (`generate-genesis.yaml`)
3. Start the bootstrap node (`start-bootstrap.yaml`)
4. Start the Lotus daemons on validator nodes (`start-daemons.yaml`)
5. Start the Mir validator processes on validator nodes (`start-validators.yaml`)

These steps are automated for convenience in the `deploy-new.yaml` playbook.
Thus, to deploy a fresh instance of Spacenet, simply run
```shell
ansible-playbook -i hosts deploy-new.yaml
```


## Provided deployment playbooks

### `clean.yaml`

Kills running Lotus daemon and Mir validator and deletes their associated state.
Does not touch the code and binaries.

Applies to all hosts by default, unless other nodes are specified using --extra-vars "nodes=..."

### `connect-daemons.yaml`

Connects all Lotus daemons to each other. This is required for the nodes to be able to sync their state.
It assumes the daemons and the bootstrap are up and running (but not necessarily the validators)

Applies to all hosts (including bootstrap) by default, unless other nodes are specified using --extra-vars "nodes=..."

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

### `generate-genesis.yaml`

Generates a new genesis block (at the bootstrap node) and copies it to all validators.
Assumes the hosts to have already been set up (using setup.yaml).

An alternative set of nodes to copy the genesis block to can be specified using --extra-vars "nodes=..."

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

Starts the bootstrap node and downloads its identity to localhost.
Assumes the host has been set up and the genesis block has been generated
(using setup.yaml and generate-genesis.yaml respectively).

Applies to the bootstrap host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-daemons.yaml`

Starts the Lotus daemons and creates connections among them and to the bootstrap node.
Assumes that the bootstrap node is up and running (see start-bootstrap.yaml).

Applies to the validator host by default, unless other nodes are specified using --extra-vars "nodes=..."

### `start-validators.yaml`

Starts the Mir validators.
Assumes that the Lotus daemons are up and running (see start-daemons.yaml).

Applies to the validator host by default, unless other nodes are specified using --extra-vars "nodes=..."
