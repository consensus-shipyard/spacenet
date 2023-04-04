package kit

import (
	"context"
	"fmt"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type FakeLotus struct {
	failedVersion bool
}

func NewFakeLotus() *FakeLotus {
	return &FakeLotus{}
}

func NewFakeLotusWithFailedVersion() *FakeLotus {
	return &FakeLotus{
		failedVersion: true,
	}
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

func (l *FakeLotus) NodeStatus(_ context.Context, _ bool) (api.NodeStatus, error) {
	return api.NodeStatus{
		SyncStatus: api.NodeSyncStatus{
			Epoch:  uint64(10),
			Behind: uint64(0),
		},
	}, nil
}

func (l *FakeLotus) Version(_ context.Context) (api.APIVersion, error) {
	if l.failedVersion {
		return api.APIVersion{}, fmt.Errorf("failed to get version")
	}
	return api.APIVersion{Version: "1.0"}, nil

}
func (l *FakeLotus) NetPeers(context.Context) ([]peer.AddrInfo, error) {
	return []peer.AddrInfo{
		{
			ID:    "ID",
			Addrs: []ma.Multiaddr{},
		},
	}, nil
}

func (l *FakeLotus) ID(context.Context) (peer.ID, error) {
	return "fakeID", nil
}
