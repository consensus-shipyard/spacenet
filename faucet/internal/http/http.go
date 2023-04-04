package http

import (
	"net/http"
	"os"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"
	"github.com/rs/cors"

	"github.com/filecoin-project/faucet/internal/faucet"
)

func FaucetHandler(logger *logging.ZapEventLogger, lotus faucet.PushWaiter, db datastore.Batching, shutdown chan os.Signal, cfg *faucet.Config) http.Handler {
	faucetService := faucet.NewService(logger, lotus, db, cfg)

	srv := NewWebService(logger, faucetService, cfg.BackendAddress)

	r := mux.NewRouter().StrictSlash(true)

	r.HandleFunc("/fund", srv.handleFunds).Methods("POST")

	r.HandleFunc("/", srv.handleHome)
	r.HandleFunc("/js/scripts.js", srv.handleScript)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	c := cors.New(cors.Options{
		AllowedOrigins:   cfg.AllowedOrigins,
		AllowCredentials: true,
	})

	return c.Handler(r)
}

func HealthHandler(logger *logging.ZapEventLogger, lotusClient LotusHealthAPI, build string, check ...ValidatorHealthCheck) http.Handler {
	h := NewHealth(logger, lotusClient, build, check...)
	r := mux.NewRouter().StrictSlash(true)
	r.HandleFunc("/readiness", h.Readiness).Methods("GET")
	r.HandleFunc("/liveness", h.Liveness).Methods("GET")
	return r
}
