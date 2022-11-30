![](./assets/spacenet-header.png)

# Spacenet
> A new generation Filecoin testnet.
>
> Made with ❤ by [ConsensusLab](https://consensuslab.world/)

- [Spacenet Faucet](https://spacenet.consensus.ninja)
- [Spacenet Genesis](./assets/genesis/spacenet.car)
- [Spacenet Bootstraps](https://github.com/consensus-shipyard/lotus/blob/spacenet/build/bootstrap/spacenet.pi)
- [Contact form](https://docs.google.com/forms/d/1O3_kHb2WJhil9sqXOxgGGGsqkAA61J1rKMfnb5os5yo/edit)

## Why Spacenet?
Spacenet is not _yet another_ Filecoin testnet. Its consensus layer has been modified to integrate [Mir](https://github.com/filecoin-project/mir), a distributed protocol implementation framework. The current version of Spacenet runs an implementation of the [Trantor BFT consensus](https://hackmd.io/P59lk4hnSBKN5ki5OblSFg) over Mir. 
_And what does this mean?_ Well, by implementing a high-performant consensus we have increased the throughput of the network significantly while reducing the 30-second block time consensus to around 1 second. As you may already be aware, Filecoin recently launched its [Filecoin Virtual Machine](https://fvm.filecoin.io/) and it will soon release support for user-defined actors. This will onboard a gamut of new applications to the Filecoin network, significantly increasing its load. Many use cases need higher throughput and tighter finality times than Filecoin can provide, and this is where Spacenet comes in.

Spacenet aims to provide developers with a testbed to deploy their FVM use cases and innovate with new Web3 applications. Many of you may be wondering, _but why would I want to develop my application over a high-throughput Filecoin network if my goal is to deploy it on mainnet?_ Well, Spacenet is just the first step towards the deployment of the [InterPlanetary Consensus (IPC)](https://github.com/filecoin-project/FIPs/discussions/419), and our ambitious plan of scaling Filecoin, both in terms of performance and new capabilities. In this very first release, Spacenet is just Filecoin with a faster consensus algorithm, but it will be much more in the future: your portal to what's to come to the Filecoin ecosystem. 

Spacenet is not only a developer sandbox to experiment with new protocol improvements, but also a way for us to test our consensus innovations with real applications and real users. Once IPC is released in Spacenet, developers will be able to deploy their own subnets from Spacenet while maintaining the ability to seamlessly interact with state and applications in the original network, from which they have otherwise become independent. With this version of Spacenet, we want to test Mir-Trantor, the first consensus algorithm supported by IPC subnets, before we unleash full IPC support. 

> In the meantime, to learn more about IPC you can read [this paper](https://research.protocol.ai/publications/hierarchical-consensus-a-horizontal-scaling-framework-for-blockchains/) and/or [watch this talk](https://www.youtube.com/watch?v=bD1LDVc2lMQ&list=PLhuBigpl7lqu0bsMQ8K7aLfmUFrkMw52K&index=3):

[![Watch the video](https://img.youtube.com/vi/bD1LDVc2lMQ/hqdefault.jpg)](https://youtu.be/bD1LDVc2lMQ)

## SLA of the network
Spacenet is an experimental network. We aim to have it constantly running, but some hiccups may appear along the way. If you are looking to rely on Spacenet for your applications, you should expect:
- Unplanned (and potentially long-lasting) downtime while we investigate bugs.
- Complete restarts of the network and loss of part or all stored state. In case of serious issues, it may be necessary to restart the network from a previous checkpoint, or completely restart from genesis.
- Bugs and rough edges to be fixed and polished along the way.

Announcements about new releases and status updates about the network are given in the #spacenet channel of the [Filecoin Slack](https://filecoin.io/slack) and through this repo. You can also ping us there or open an issue in this repo if you encounter a bug or some other issue with the network. You can also direct your requests through [this form](https://docs.google.com/forms/d/1O3_kHb2WJhil9sqXOxgGGGsqkAA61J1rKMfnb5os5yo/edit).

## Getting started for users
Spacenet is a Filecoin testnet, and as such it is supposed to do (almost) everything that the [Filecoin network supports](https://lotus.filecoin.io/tutorials/lotus/store-and-retrieve/set-up/):
- Send Filecoin between addresses.
- Create multisig accounts.
- Create [payment channels](https://lotus.filecoin.io/tutorials/lotus/payment-channels/).
- Deploy [FVM contracts](https://docs.filecoin.io/fvm/basics/introduction/).

That being said, as the consensus layer is no longer storage-dependent, Spacenet has limited support for storage-related features. In particular, we have stripped out some of the functionalities of the lotus miner. While you deploy a lotus-miner over Spacenet to onboard storage to the network and perform deals, lotus-miners are not allowed to propose and validate blocks anymore (this is handled by Mir-Trantor validators).

> ⚠️ Support for storage-specific features in Spacenet is limited.

### Getting Spacenet FIL
In order to fund your account with Spacenet FIL we provide a faucet at [https://spacenet.consensus.ninja](https://spacenet.consensus.ninja). Getting FIL is as simple as inputting your address in the textbox and clicking the button.
- The per-request allowance given by the faucet is of 10 FIL.
- There is a daily maximum of 20 FIL per address.
- And we have also limited the maximum amount of funds that the faucet can withdraw daily.
If, for some reason, you require more Spacenet FIL for your application, feel free to drop us a message in the #spacenet Slack channel, via consensuslab@protocol.ai to increase your allowance, or fill-in a request in [this form](https://docs.google.com/forms/d/1O3_kHb2WJhil9sqXOxgGGGsqkAA61J1rKMfnb5os5yo/edit).
![](./assets/spacenet-faucet.png)

## Getting started for developers
You can run a full-node and connect it to Spacenet by:
- Cloning the modified lotus implementation for Spacenet:
```
git clone https://github.com/consensus-shipyard/lotus

// The latest stable branch for the network is `spacenet`
git checkout spacenet
```
- Installing lotus and running all dependencies as described in the `README` of the [repo](https://github.com/consensus-shipyard/lotus)
- Once you have all `lotus` installed you can run the following command to compile `lotus` with Spacenet support.
```
make spacenet
```
- With that, you are ready to run your spacenet daemon and connect to the network by connecting to any its bootstrap nodes.
```
./lotus daemon --bootsraps=true
```
Spacenet supports every lotus command supported in mainnet, so you'll be able to configure your Spacenet full-node at will (by exposing a different API port, running Lotus lite, etc.). More info available in [Lotus' docs](https://lotus.filecoin.io/lotus/get-started/what-is-lotus/).

## Getting started for validators

> Support for external validators coming soon! Track the work in [the following issue](https://github.com/consensus-shipyard/lotus/issues/21). If you are interested in becoming a validator let us know through [this form](https://docs.google.com/forms/d/1O3_kHb2WJhil9sqXOxgGGGsqkAA61J1rKMfnb5os5yo).

Spacenet is currently run by a committee of 4 validators owned by CL. We don't accept externally owned validators during this initial testing phase, until the network deployment is stabilized, but support for reconfiguration and external validators will be added soon.

## What's next?
- Reconfiguration.
- FEVM.
- Native WASM actors.


