package kit

import (
	"context"

	"github.com/ipfs/go-cid"

	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/chain/types"
)

type FakeLotus struct {
}

func NewFakeLotus() *FakeLotus {
	return &FakeLotus{}
}

func (l *FakeLotus) MpoolPushMessage(ctx context.Context, msg *types.Message, spec *api.MessageSendSpec) (*types.SignedMessage, error) {
	smsg := types.SignedMessage{
		Message: *msg,
	}
	return &smsg, nil
}
func (l *FakeLotus) StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*api.MsgLookup, error) {
	return nil, nil
}
