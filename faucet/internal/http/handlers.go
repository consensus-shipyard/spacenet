package http

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/db"
	"github.com/filecoin-project/faucet/internal/platform/web"
	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
)

var (
	TotalMax    = uint64(4000)
	AddressMax  = uint64(2000)
	TokenAmount = abi.NewTokenAmount(1000)
)

type FaucetService struct {
	log    *logging.ZapEventLogger
	lotus  v0api.FullNode
	db     *db.Database
	faucet address.Address
}

func NewFaucetService(log *logging.ZapEventLogger, lotus v0api.FullNode, store datastore.Datastore, faucet address.Address) *FaucetService {
	return &FaucetService{
		log:    log,
		lotus:  lotus,
		db:     db.NewDatabase(store),
		faucet: faucet,
	}
}

func (h *FaucetService) fundable(ctx context.Context, targetAddr address.Address) error {
	addrInfo, err := h.db.GetAddrInfo(ctx, targetAddr)
	if err != nil {
		return err
	}
	h.log.Infof("target address info: %v", addrInfo)

	totalInfo, err := h.db.GetTotalInfo(ctx)
	if err != nil {
		return err
	}
	h.log.Infof("total info: %v", totalInfo)

	if time.Since(addrInfo.LatestWithdrawal) >= 24*time.Hour {
		addrInfo.Amount = 0
		addrInfo.LatestWithdrawal = time.Now()
	}

	if time.Since(totalInfo.LatestWithdrawal) >= 24*time.Hour {
		totalInfo.Amount = 0
		totalInfo.LatestWithdrawal = time.Now()
	}

	if totalInfo.Amount < TotalMax && addrInfo.Amount < AddressMax {
		h.log.Infof("transaction is allowed")
		addrInfo.Amount += TokenAmount.Uint64()
		totalInfo.Amount += TokenAmount.Uint64()
	} else {
		return fmt.Errorf("transaction exceeds allowed funds")
	}

	if err = h.db.UpdateAddrInfo(ctx, targetAddr, addrInfo); err != nil {
		return err
	}

	if err = h.db.UpdateTotalInfo(ctx, totalInfo); err != nil {
		return err
	}
	h.log.Infof("total info: %v", totalInfo)

	return nil
}

func (h *FaucetService) fund(w http.ResponseWriter, r *http.Request) {
	var req data.FundRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	h.log.Infof("input request: %v", req)

	targetAddr, err := address.NewFromString(req.Address)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = h.fundable(r.Context(), targetAddr)
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	err = h.pushMessage(r.Context(), targetAddr)
	if err != nil {
		h.log.Errorw("Error waiting for message to be committed", "err", err)
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	return

}

func (h *FaucetService) pushMessage(ctx context.Context, addr address.Address) error {
	msg, err := h.lotus.MpoolPushMessage(ctx, &types.Message{
		To:     addr,
		From:   h.faucet,
		Value:  TokenAmount,
		Method: 0, // method Send
		Params: nil,
	}, nil)
	if err != nil {
		return err
	}

	// wait state message.
	if _, err = h.lotus.StateWaitMsg(ctx, msg.Cid(), build.MessageConfidence); err != nil {
		return err
	}

	h.log.Infow("Address funded successfully")
	return nil
}
