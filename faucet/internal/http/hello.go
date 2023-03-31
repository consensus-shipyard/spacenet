package http

import (
	"context"
	"net/http"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p-core/peer"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/platform/web"
	"github.com/filecoin-project/lotus/api"
)

type Status interface {
	NodeStatus(ctx context.Context, inclChainStatus bool) (api.NodeStatus, error)
	NetPeers(context.Context) ([]peer.AddrInfo, error)
}

type HelloService struct {
	log *logging.ZapEventLogger
	srv Status
}

func NewHelloService(log *logging.ZapEventLogger, srv Status) *HelloService {
	return &HelloService{
		log: log,
		srv: srv,
	}
}

func (h *HelloService) handleHello(w http.ResponseWriter, r *http.Request) {

	h.log.Infof(">>> %s -> {%s}\n", r.RemoteAddr, r.URL)

	status, err := h.srv.NodeStatus(r.Context(), true)
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	p, err := h.srv.NetPeers(r.Context())
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	resp := data.HelloResponse{
		Epoch:      status.SyncStatus.Epoch,
		PeerNumber: len(p),
	}

	if err := web.Respond(r.Context(), w, resp, http.StatusOK); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
}
