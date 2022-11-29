package http

import (
	"errors"
	"html/template"
	"net/http"
	"path"

	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/faucet"
	"github.com/filecoin-project/faucet/internal/platform/web"
	"github.com/filecoin-project/go-address"
)

type WebService struct {
	log            *logging.ZapEventLogger
	faucet         *faucet.Service
	backendAddress string
}

func NewWebService(log *logging.ZapEventLogger, faucet *faucet.Service, backendAddress string) *WebService {
	return &WebService{
		log:            log,
		faucet:         faucet,
		backendAddress: backendAddress,
	}
}

func (h *WebService) handleFunds(w http.ResponseWriter, r *http.Request) {
	var req data.FundRequest
	if err := web.Decode(r, &req); err != nil {
		web.RespondError(w, http.StatusBadRequest, err)
		return
	}

	if req.Address == "" {
		web.RespondError(w, http.StatusBadRequest, errors.New("empty address"))
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
}

func (h *WebService) handleHome(w http.ResponseWriter, r *http.Request) {
	p := path.Dir("./static/index.html")
	w.Header().Set("Content-type", "text/html")
	http.ServeFile(w, r, p)
}

func (h *WebService) handleScript(w http.ResponseWriter, r *http.Request) {
	tmpl, err := template.ParseFiles("./static/js/scripts.js")
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	if err = tmpl.Execute(w, h.backendAddress); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
	w.Header().Set("Content-type", "text/javascript")
}
