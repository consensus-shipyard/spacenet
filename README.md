![](./assets/spacenet-header.png)

# Spacenet
> A new generation Filecoin testnet.
>
> Made with ❤ by [ConsensusLab](https://consensuslab.world/)

- Spacenet Faucet: https://consensus-shipyard.github.io/spacenet-faucet/
- [Spacenet Genesis](./assets/genesis/spacenet.car)
- Spacenet Bootstraps:
  - `xxxx`
  - `xxxx`

## Why Spacenet?
Spacenet is not _yet another_ Filecoin testnet. Its consensus layer has been modified to integrate [Mir](https://github.com/filecoin-project/mir), a distributed protocol implementation framework. The current version of Spacenet runs an implementation of the Trantor BFT consensus over Mir. 
_And what does this mean?_ Well, by implementing a high-performant consensus we have increased the throughput of the network significantly while reducing the 30 second block time consensus to around 1 second. As you may already be aware, Filecoin recently launched its [Filecoin Virtual Machine](https://fvm.filecoin.io/) and it will soon release support for user-defined actor. This will on board a great gamut of new applications to the Filecoin network significantly increasing its load. Certain use cases may realize that they need more from Filecoin in terms of throughput and finality times, and here is were Spacenet comes in.

With Spacenet we want to provide developers with a testbed to deploy their FVM use cases and innovate with new Web3 applications. Many of you may be wondering, _but why would I want to develop my application over a high-throughput Filecoin network if my goal is to deploy it on mainnet?_ Well, Spacenet is just the first step towards the deployment of the [InterPlanetary Consensus (IPC)](https://github.com/filecoin-project/FIPs/discussions/419), and our ambitious plan of scaling Filecoin (in terms of performance and new capabilities). Now Spacenet is just Filecoin with a faster consensus algorithm, but in the future it will be much more. Spacenet is just your portal to glimpse what is yet to come to the Filecoin ecosystem (but we are getting ahead of ourselves). 

Spacenet is not only a testbed for developer to test the new improvements to the protocol we are working on, it is also a way for us to test our consensus innovations with real applications and real users. For instance, once IPC is released in Spacenet, developers will be able to deploy their own subnets from Spacenet while maintaining the ability to seamlessly interact with state and applications in the original network from which they have become independent. With this version of Spacenet, we want to test Mir-Trantor, the first consensus algorithm supported by IPC subnets, before we release full support for IPC. 

> In the meantime, to learn more about IPC you can read [this paper](https://research.protocol.ai/publications/hierarchical-consensus-a-horizontal-scaling-framework-for-blockchains/) and/or [watch this talk](https://www.youtube.com/watch?v=bD1LDVc2lMQ&list=PLhuBigpl7lqu0bsMQ8K7aLfmUFrkMw52K&index=3).

## SLA of the network
Spacenet is an experimental network. We'll do our best to have it always running, but some hiccups may appear along the way. If you are looking to rely on Spacenet for your applications you should expect:
- Unexpected (and potentially long-lasting) downtimes while we investigate bugs.
- Complete restarts of the network and loss of part or all stored state. If really nasty bugs appear we may need to restart the network from a previous checkpoint, or completely restart from genesis.
- Bugs and rough edges to be fixed and polished along the way.

Announcements about new releases and status updates about the network are given in the #spacenet channel of the [Filecoin Slack](filecoinproject.slack.com) and through this repo. You can also ping us there or open an issue in this repo if you encounter a bug or some other issue with the network.

## Getting started for users
Spacenet is a Filecoin testnet, and as such it is supposed to do (almost) everything that a [Filecoin network support](https://lotus.filecoin.io/tutorials/lotus/store-and-retrieve/set-up/).
- Send Filecoin between addresses.
- Create multisig accounts.
- Create [payment channels](https://lotus.filecoin.io/tutorials/lotus/payment-channels/).
- Deploy [FVM contracts](https://docs.filecoin.io/fvm/basics/introduction/).

That being said, as the consensus layer is no longer storage-dependent, Spacenet has limited support for storage-related features. As the consensus algorithm is no longer storage-dependent, we have stripped-out some of the functionalities of the lotus miner. You can deploy a lotus-miner over Spacenet to on board storage to the network and perform deals. However, lotus-miners are not allowed to propose and validate blocks anymore (this is handled by Mir-Trantor validators).

> ⚠️ Support for storage-specific features in Spacenet is limited.

### Getting Spacenet FIL
In order to fund your account with Spacenet FIL we provide a faucet at [https://consensus-shipyard.github.io/spacenet-faucet/](https://consensus-shipyard.github.io/spacenet-faucet/). Getting FIL is as simple inputting your address in the textbox and clicking the button.
- The per-request allowance given by the faucet is of 10FIL.
- There is a daily maximum of 20FIL per address.
- And we have also limited the maximum amount of funds that the faucet can withdraw daily.
If for some reason you require more Spacenet FIL for your application, feel free to drop us a message at consensuslab@protocol.ai or the #spacenet Slack channel to increase your allowance.
![](./assets/spacenet-faucet.png)

## Getting started for developers
You can run a full-node and connect it to Spacenet by:
- Cloning the modified lotus implementation for Spacenet:
```
git clone https://github.com/consensus-shipyard/lotus
// TODO: Remove this step once we merge our version to master.
git checkout adlrocha/mir-sync
```
- Installing lotus and running all dependencies as described in the `README` of the [repo](https://github.com/consensus-shipyard/lotus)
- Once you have all `lotus` installed you can run the following command to compile `lotus` with Spacenet support.
```
make spacenet
```
- With that, you are ready to run your spacenet daemon and connect to the network by connecting to any its bootstrapper nodes.
```
./lotus daemon --genesis=spacenet.car
```

> TODO: Point to spacenet bootstraps.

## Getting started for validators

> Support for external validators coming soon!

Spacenet is currently run by a committee of 4 validators owned by CL. Initially, we don't accept externally owned validators until the network deployment is stabilized, but support for reconfiguration and externally will be released soon.

## What's next?
- Support for reconfiguration.
- FEVM support.
- Native WASM actors support.


