package http

import (
	"net/http"
	"path"

	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/faucet"
	"github.com/filecoin-project/faucet/internal/platform/web"
	"github.com/filecoin-project/go-address"
)

type WebService struct {
	log    *logging.ZapEventLogger
	faucet *faucet.Service
}

func NewWebService(log *logging.ZapEventLogger, faucet *faucet.Service) *WebService {
	return &WebService{
		log:    log,
		faucet: faucet,
	}
}

func (h *WebService) handleFunds(w http.ResponseWriter, r *http.Request) {
	var req data.FundRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	h.log.Infof("Input request: %v", req)

	targetAddr, err := address.NewFromString(req.Address)
	if err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	err = h.faucet.FundAddress(r.Context(), targetAddr)
	if err != nil {
		h.log.Errorw("Failed to fund address", "addr", targetAddr, "err", err)
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	w.WriteHeader(http.StatusCreated)
	w.Header().Set("Content-Type", "application/json")
}

func (h *WebService) handleHome(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./static/index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}
