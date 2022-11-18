package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/faucet"
	"github.com/filecoin-project/go-address"
)

func Handler(log *logging.ZapEventLogger, lotus faucet.PushWaiter, db datastore.Batching, shutdown chan os.Signal, faucetAddr address.Address) http.Handler {
	faucetService := faucet.NewService(log, lotus, db, faucetAddr)

	srv := NewWebService(log, faucetService)

	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/fund", srv.handleFunds).Methods("POST")
	r.HandleFunc("/", srv.handleHome)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	return r
}
