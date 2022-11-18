package faucet

import (
	"context"
	"fmt"
	"time"

	"github.com/ipfs/go-cid"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/db"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
)

var (
	// Amounts of tokens in FIL.
	TotalWithdrawalLimit   = uint64(1000)
	AddressWithdrawalLimit = uint64(20)
	WithdrawalAmount       = uint64(10)

	ErrExceedTotalAllowed = fmt.Errorf("transaction exceeds total allowed funds per day")
	ErrExceedAddrAllowed  = fmt.Errorf("transaction to exceeds daily allowed funds per address")
)

type Service struct {
	log    *logging.ZapEventLogger
	lotus  PushWaiter
	db     *db.Database
	faucet address.Address
}

type PushWaiter interface {
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *api.MessageSendSpec) (*types.SignedMessage, error)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*api.MsgLookup, error)
}

func NewService(log *logging.ZapEventLogger, lotus PushWaiter, store datastore.Datastore, faucet address.Address) *Service {
	return &Service{
		log:    log,
		lotus:  lotus,
		db:     db.NewDatabase(store),
		faucet: faucet,
	}
}

func (s *Service) FundAddress(ctx context.Context, targetAddr address.Address) error {
	addrInfo, err := s.db.GetAddrInfo(ctx, targetAddr)
	if err != nil {
		return err
	}
	s.log.Infof("target address info: %v", addrInfo)

	totalInfo, err := s.db.GetTotalInfo(ctx)
	if err != nil {
		return err
	}
	s.log.Infof("total info: %v", totalInfo)

	if time.Since(addrInfo.LatestWithdrawal) >= 24*time.Hour {
		addrInfo.Amount = 0
		addrInfo.LatestWithdrawal = time.Now()
	}

	if time.Since(totalInfo.LatestWithdrawal) >= 24*time.Hour {
		totalInfo.Amount = 0
		totalInfo.LatestWithdrawal = time.Now()
	}

	if totalInfo.Amount >= TotalWithdrawalLimit {
		return ErrExceedTotalAllowed
	}

	if addrInfo.Amount >= AddressWithdrawalLimit {
		return ErrExceedAddrAllowed
	}

	s.log.Infof("funding %v is allowed", targetAddr)

	err = s.pushMessage(ctx, targetAddr)
	if err != nil {
		s.log.Errorw("Error waiting for message to be committed", "err", err)
		return fmt.Errorf("failt to push message: %w", err)
	}

	addrInfo.Amount += WithdrawalAmount
	totalInfo.Amount += WithdrawalAmount

	if err = s.db.UpdateAddrInfo(ctx, targetAddr, addrInfo); err != nil {
		return err
	}

	if err = s.db.UpdateTotalInfo(ctx, totalInfo); err != nil {
		return err
	}

	return nil
}

func (s *Service) pushMessage(ctx context.Context, addr address.Address) error {
	msg, err := s.lotus.MpoolPushMessage(ctx, &types.Message{
		To:     addr,
		From:   s.faucet,
		Value:  types.FromFil(WithdrawalAmount),
		Method: 0, // method Send
		Params: nil,
	}, nil)
	if err != nil {
		return err
	}

	if _, err = s.lotus.StateWaitMsg(ctx, msg.Cid(), build.MessageConfidence); err != nil {
		return err
	}

	s.log.Infof("Address %v funded successfully", addr)
	return nil
}
