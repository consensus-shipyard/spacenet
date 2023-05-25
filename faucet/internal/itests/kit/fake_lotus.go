package kit

import (
	"context"
	"fmt"
	"sync"

	"github.com/ipfs/go-cid"
	"github.com/libp2p/go-libp2p/core/peer"
	ma "github.com/multiformats/go-multiaddr"

	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type FakeLotus struct {
	m             sync.Mutex
	failedVersion bool
	h             uint64
	failed        bool
	failedOn      uint64
}

func NewFakeLotus(failed bool, failedOn uint64) *FakeLotus {
	return &FakeLotus{
		failed:   failed,
		failedOn: failedOn,
	}
}

func NewFakeLotusNoCrash() *FakeLotus {
	return &FakeLotus{
		failed:   false,
		failedOn: 0,
	}
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
	l.m.Lock()
	defer l.m.Unlock()

	s := api.NodeStatus{
		SyncStatus: api.NodeSyncStatus{
			Epoch:  l.h,
			Behind: uint64(0),
		},
	}
	if !l.failed {
		l.h++
	} else {
		if l.h < l.failedOn {
			l.h++
		}
	}
	return s, nil
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
