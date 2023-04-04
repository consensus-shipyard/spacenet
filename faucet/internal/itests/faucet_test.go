package itests

import (
	"bytes"
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"os"
	"strings"
	"testing"
	"time"

	datastore "github.com/ipfs/go-ds-leveldb"
	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"
	ldbopts "github.com/syndtr/goleveldb/leveldb/opt"

	"github.com/filecoin-project/faucet/internal/data"
	faucetDB "github.com/filecoin-project/faucet/internal/db"
	"github.com/filecoin-project/faucet/internal/faucet"
	handler "github.com/filecoin-project/faucet/internal/http"
	"github.com/filecoin-project/faucet/internal/itests/kit"
	"github.com/filecoin-project/go-address"
)

type FaucetTests struct {
	handler   http.Handler
	store     *datastore.Datastore
	db        *faucetDB.Database
	faucetCfg *faucet.Config
}

const (
	FaucetAddr = "f1cp4q4lqsdhob23ysywffg2tvbmar5cshia4rweq"
	TestAddr1  = "f1akaouty2buxxwb46l27pzrhl3te2lw5jem67xuy"
	TestAddr2  = "f1vfp7yzvwy7ftktnex2cfoz2gpm2jyxlebqpam4q"
	storePath  = "./_store"
)

func Test_Faucet(t *testing.T) {
	store, err := datastore.NewDatastore(storePath, &datastore.Options{
		Compression: ldbopts.NoCompression,
		NoSync:      false,
		Strict:      ldbopts.StrictAll,
		ReadOnly:    false,
	})
	require.NoError(t, err)

	defer func() {
		err = store.Close()
		require.NoError(t, err)
		err = os.RemoveAll(storePath)
		require.NoError(t, err)
	}()

	log := logging.Logger("TEST-FAUCET")

	lotus := kit.NewFakeLotus()

	addr, err := address.NewFromString(FaucetAddr)
	require.NoError(t, err)

	shutdown := make(chan os.Signal, 1)

	cfg := faucet.Config{
		FaucetAddress:          addr,
		TotalWithdrawalLimit:   1000,
		AddressWithdrawalLimit: 20,
		WithdrawalAmount:       10,
	}

	srv := handler.FaucetHandler(log, lotus, store, shutdown, &cfg)

	db := faucetDB.NewDatabase(store)

	tests := FaucetTests{
		handler:   srv,
		store:     store,
		db:        db,
		faucetCfg: &cfg,
	}

	t.Run("fundEmptyAddress", tests.emptyAddress)
	t.Run("fundAddress201", tests.fundAddress201)
	t.Run("fundAddressWithMoreThanAllowed", tests.fundAddressWithMoreThanAllowed)
	t.Run("fundAddressWithMoreThanTotal", tests.fundAddressWithMoreThanTotal)
}

func (ft *FaucetTests) emptyAddress(t *testing.T) {
	req := data.FundRequest{Address: ""}

	body, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/fund", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ft.handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusBadRequest, w.Code)
}

func (ft *FaucetTests) fundAddress201(t *testing.T) {
	req := data.FundRequest{Address: FaucetAddr}

	body, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/fund", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ft.handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusCreated, w.Code)
}

// fundAddressWithMoreThanAllowed tests that exceeding daily allowed funds per address is not allowed.
func (ft *FaucetTests) fundAddressWithMoreThanAllowed(t *testing.T) {
	targetAddr, err := address.NewFromString(TestAddr1)
	require.NoError(t, err)

	err = ft.db.UpdateAddrInfo(context.Background(), targetAddr, data.AddrInfo{
		Amount:           ft.faucetCfg.AddressWithdrawalLimit,
		LatestWithdrawal: time.Now(),
	})
	require.NoError(t, err)

	req := data.FundRequest{Address: TestAddr1}

	body, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/fund", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ft.handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	got := w.Body.String()
	exp := faucet.ErrExceedAddrAllowedFunds.Error()
	if !strings.Contains(got, exp) {
		t.Logf("\t\tTest %s:\tGot : %v", t.Name(), got)
		t.Logf("\t\tTest %s:\tExp: %v", t.Name(), exp)
		t.Fatalf("\t\tTest %s:\tShould get the expected result.", t.Name())
	}
}

// fundAddressWithMoreThanAllowed tests that exceeding daily allowed funds per address is not allowed.
func (ft *FaucetTests) fundAddressWithMoreThanTotal(t *testing.T) {
	err := ft.db.UpdateTotalInfo(context.Background(), data.TotalInfo{
		Amount:           ft.faucetCfg.TotalWithdrawalLimit,
		LatestWithdrawal: time.Now(),
	})
	require.NoError(t, err)

	req := data.FundRequest{Address: TestAddr2}

	body, err := json.Marshal(&req)
	if err != nil {
		t.Fatal(err)
	}

	r := httptest.NewRequest(http.MethodPost, "/fund", bytes.NewBuffer(body))
	w := httptest.NewRecorder()

	ft.handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusInternalServerError, w.Code)

	got := w.Body.String()
	exp := faucet.ErrExceedTotalAllowedFunds.Error()
	if !strings.Contains(got, exp) {
		t.Logf("\t\tTest %s:\tGot : %v", t.Name(), got)
		t.Logf("\t\tTest %s:\tExp: %v", t.Name(), exp)
		t.Fatalf("\t\tTest %s:\tShould get the expected result.", t.Name())
	}
}
