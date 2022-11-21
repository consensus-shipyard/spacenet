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
	ErrExceedTotalAllowedFunds = fmt.Errorf("transaction exceeds total allowed funds per day")
	ErrExceedAddrAllowedFunds  = fmt.Errorf("transaction to exceeds daily allowed funds per address")
)

type PushWaiter interface {
	MpoolPushMessage(ctx context.Context, msg *types.Message, spec *api.MessageSendSpec) (*types.SignedMessage, error)
	StateWaitMsg(ctx context.Context, cid cid.Cid, confidence uint64) (*api.MsgLookup, error)
}

type Config struct {
	FaucetAddress          address.Address
	TotalWithdrawalLimit   uint64
	AddressWithdrawalLimit uint64
	WithdrawalAmount       uint64
}

type Service struct {
	log    *logging.ZapEventLogger
	lotus  PushWaiter
	db     *db.Database
	faucet address.Address
	cfg    *Config
}

func NewService(log *logging.ZapEventLogger, lotus PushWaiter, store datastore.Datastore, cfg *Config) *Service {
	return &Service{
		cfg:    cfg,
		log:    log,
		lotus:  lotus,
		db:     db.NewDatabase(store),
		faucet: cfg.FaucetAddress,
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

	if addrInfo.LatestWithdrawal.IsZero() || time.Since(addrInfo.LatestWithdrawal) >= 24*time.Hour {
		addrInfo.Amount = 0
		addrInfo.LatestWithdrawal = time.Now()
	}

	if totalInfo.LatestWithdrawal.IsZero() || time.Since(totalInfo.LatestWithdrawal) >= 24*time.Hour {
		totalInfo.Amount = 0
		totalInfo.LatestWithdrawal = time.Now()
	}

	if totalInfo.Amount >= s.cfg.TotalWithdrawalLimit {
		return ErrExceedTotalAllowedFunds
	}

	if addrInfo.Amount >= s.cfg.AddressWithdrawalLimit {
		return ErrExceedAddrAllowedFunds
	}

	s.log.Infof("funding %v is allowed", targetAddr)

	err = s.pushMessage(ctx, targetAddr)
	if err != nil {
		s.log.Errorw("Error waiting for message to be committed", "err", err)
		return fmt.Errorf("failt to push message: %w", err)
	}

	addrInfo.Amount += s.cfg.WithdrawalAmount
	totalInfo.Amount += s.cfg.WithdrawalAmount

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
		Value:  types.FromFil(s.cfg.WithdrawalAmount),
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
