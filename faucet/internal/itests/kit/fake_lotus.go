package kit

import (
	"context"

	"github.com/ipfs/go-cid"

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
