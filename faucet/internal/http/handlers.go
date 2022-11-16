package http

import (
	"context"
	"encoding/binary"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"time"

	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/go-state-types/abi"
	"github.com/filecoin-project/lotus/api/v0api"
	"github.com/filecoin-project/lotus/build"
	"github.com/filecoin-project/lotus/chain/types"
)

var (
	TotalMax                 = uint64(4000)
	AddressMax               = uint64(2000)
	TokenAmount              = abi.NewTokenAmount(1000)
	TotalAmountKey           = datastore.NewKey("totalAmount")
	TotalLatestWithdrawalKey = datastore.NewKey("totalLatestWithdrawal")
)

type SpaceService struct {
	log    *logging.ZapEventLogger
	lotus  v0api.FullNode
	db     datastore.Datastore
	faucet address.Address
}

func NewSpaceService(log *logging.ZapEventLogger, lotus v0api.FullNode, db datastore.Datastore, faucet address.Address) *SpaceService {
	return &SpaceService{
		log:    log,
		lotus:  lotus,
		db:     db,
		faucet: faucet,
	}
}

type fundRequest struct {
	// Value abi.TokenAmount
	Address string
}

type fundResponse struct {
	Error string
}

func (h *SpaceService) fundable(targetAddr address.Address, targetValue abi.TokenAmount) error {
	var totalAmount uint64
	b, err := h.db.Get(context.TODO(), TotalAmountKey)
	if err != nil && err != datastore.ErrNotFound {
		return fmt.Errorf("failed to get total amount: %w", err)
	}
	if err == datastore.ErrNotFound {
		totalAmount = 0
	} else {
		totalAmount = binary.BigEndian.Uint64(b)
	}
	h.log.Infof("total amount: %v", totalAmount)

	var targetAmount uint64
	b, err = h.db.Get(context.TODO(), AddrKey(targetAddr))
	if err != nil && err != datastore.ErrNotFound {
		return fmt.Errorf("failed to get addr token value: %w", err)
	}
	if err == datastore.ErrNotFound {
		targetAmount = 0
	} else {
		targetAmount = binary.BigEndian.Uint64(b)
	}

	h.log.Infof("%v address amount: %v", targetAddr, targetAmount)

	var targetLatestWithdrawal time.Time
	b, err = h.db.Get(context.TODO(), LatestWithdrawalKey(targetAddr))
	if err != nil && err != datastore.ErrNotFound {
		return fmt.Errorf("failed to get latest withdrawal: %w", err)
	}
	if err == datastore.ErrNotFound {
		targetLatestWithdrawal = time.Now().Add(-time.Hour * 24)
	} else {
		err = targetLatestWithdrawal.UnmarshalBinary(b)
		if err != nil && err != datastore.ErrNotFound {
			return fmt.Errorf("failed to unmarshal latest withdrawal: %w", err)
		}
	}

	h.log.Infof("%v address latest withdrawal: %v", targetAddr, targetLatestWithdrawal)

	var totalLatestWithdrawal time.Time
	b, err = h.db.Get(context.TODO(), TotalLatestWithdrawalKey)
	if err != nil && err != datastore.ErrNotFound {
		return fmt.Errorf("failed to get total latest withdrawal: %v", err)
	}
	if err == datastore.ErrNotFound {
		totalLatestWithdrawal = time.Now().Add(-time.Hour * 24)
	} else {
		err = totalLatestWithdrawal.UnmarshalBinary(b)
		if err != nil && err != datastore.ErrNotFound {
			return fmt.Errorf("failed to unmarshal total latest withdrawal: %w", err)
		}
	}

	h.log.Infof("latest total withdrawal: %v", totalLatestWithdrawal)

	if time.Since(targetLatestWithdrawal) >= 24*time.Hour {
		fmt.Println(1)
		targetAmount = 0
		targetLatestWithdrawal = time.Now()
	}

	if time.Since(totalLatestWithdrawal) >= 24*time.Hour {
		fmt.Println(2)
		totalAmount = 0
		totalLatestWithdrawal = time.Now()
	}

	if totalAmount < TotalMax && targetAmount < AddressMax {
		fmt.Println(3)
		totalAmount += targetValue.Uint64()
		targetAmount += targetValue.Uint64()
	} else {
		fmt.Println(5)
		return fmt.Errorf("transaction exceeds allowed funds")
	}

	fmt.Println(4)

	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, targetAmount)
	err = h.db.Put(context.TODO(), AddrKey(targetAddr), b)
	if err != nil {
		return fmt.Errorf("failed to put target amount: %v", err)
	}

	b = make([]byte, 8)
	binary.LittleEndian.PutUint64(b, totalAmount)
	err = h.db.Put(context.TODO(), TotalAmountKey, b)
	if err != nil {
		return fmt.Errorf("failed to put total amount: %v", err)
	}

	b, err = targetLatestWithdrawal.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal latest withdrawal")
	}
	err = h.db.Put(context.TODO(), LatestWithdrawalKey(targetAddr), b)
	if err != nil {
		return fmt.Errorf("failed to put latest withdrawal: %v", err)
	}

	b, err = totalLatestWithdrawal.MarshalBinary()
	if err != nil {
		return fmt.Errorf("failed to marshal total latest withdrawal")
	}
	err = h.db.Put(context.TODO(), TotalLatestWithdrawalKey, b)
	if err != nil {
		return fmt.Errorf("failed to put latest total withdrawal: %v", err)
	}

	return nil
}

// TODO: Finalize this method.
func (h *SpaceService) fundRequest(w http.ResponseWriter, r *http.Request) {
	body, err := io.ReadAll(r.Body)
	if err != nil {
		h.log.Errorf("error decoding request: %w", err)
		return
	}
	fmt.Println(">>>>>>>> Body", string(body))
	var req fundRequest
	err = json.Unmarshal(body, &req)
	if err != nil {
		h.log.Errorf("error unmarshaling json request: %w", err)
		errResponse(w, "error unmarshaling request")
		return
	}

	addr, err := address.NewFromString(req.Address)
	if err != nil {
		errResponse(w, err.Error())
		return
	}
	err = h.fundable(addr, TokenAmount)
	if err != nil {
		errResponse(w, err.Error())
		return
	}
	err = h.fundAddr(addr, TokenAmount)
	if err != nil {
		errResponse(w, err.Error())
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
	return

}

func errResponse(w http.ResponseWriter, errStr string) {
	resp := fundResponse{Error: errStr}
	jsonResp, err := json.Marshal(resp)
	if err != nil {
		log.Fatalf("Error happened in JSON marshal. Err: %s", err)
	}
	w.Write(jsonResp)
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
}

func (h *SpaceService) fundAddr(addr address.Address, value abi.TokenAmount) error {
	ctx := context.TODO()

	msg, err := h.lotus.MpoolPushMessage(ctx, &types.Message{
		To:     addr,
		From:   h.faucet,
		Value:  value,
		Method: 0, // methodSend
		Params: nil,
	}, nil)
	if err != nil {
		h.log.Errorw("Error pushing join subnet message to parent api", "err", err)
		return err
	}

	// wait state message.
	_, err = h.lotus.StateWaitMsg(ctx, msg.Cid(), build.MessageConfidence)
	if err != nil {
		h.log.Errorw("Error waiting for message to be committed", "err", err)
		return err
	}

	h.log.Infow("Address funded successfully")
	return nil
}

func AddrKey(addr address.Address) datastore.Key {
	return datastore.NewKey(addr.String() + ":value")
}

func LatestWithdrawalKey(addr address.Address) datastore.Key {
	return datastore.NewKey(addr.String() + ":latestWithdrawal")
}
