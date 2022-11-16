package http

import (
	"net/http"
	"os"
	"path"

	"github.com/gorilla/mux"
	"github.com/ipfs/go-datastore"
	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/lotus/api/v0api"
)

func Handler(log *logging.ZapEventLogger, lotus v0api.FullNode, db datastore.Batching, shutdown chan os.Signal, faucet address.Address) http.Handler {
	r := mux.NewRouter().StrictSlash(true)

	srv := NewSpaceService(log, lotus, db, faucet)

	r.HandleFunc("/fund", srv.fundRequest).Methods("POST")
	r.HandleFunc("/", home)
	r.PathPrefix("/").Handler(http.StripPrefix("/", http.FileServer(http.Dir("./static"))))

	return r
}

// serves index file
func home(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./static/index.html")
	// set header
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}
