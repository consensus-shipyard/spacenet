package kit

import (
	"context"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p-core/peer"
	"github.com/libp2p/go-libp2p/core"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type FakeLotus struct {
}

func NewFakeLotus() *FakeLotus {
	return &FakeLotus{}
}

func (l *FakeLotus) MpoolPushMessage(_ context.Context, msg *types.Message, _ *api.MessageSendSpec) (*types.SignedMessage, error) {
	smsg := types.SignedMessage{
		Message: *msg,
	}
	return &smsg, nil
}

func (l *FakeLotus) StateWaitMsg(_ context.Context, _ cid.Cid, _ uint64, _ abi.ChainEpoch, _ bool) (*api.MsgLookup, error) {
	return nil, nil
}

func (l *FakeLotus) NodeStatus(ctx context.Context, inclChainStatus bool) (api.NodeStatus, error) {
	return api.NodeStatus{
		SyncStatus: api.NodeSyncStatus{
			Epoch:  uint64(10),
			Behind: uint64(0),
		},
	}, nil

}
func (l *FakeLotus) NetPeers(context.Context) ([]peer.AddrInfo, error) {
	return []peer.AddrInfo{
		{
			ID:    "ID",
			Addrs: []core.Multiaddr{},
		},
	}, nil
}
