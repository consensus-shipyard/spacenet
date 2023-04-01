package http

import (
	"context"
	"net/http"
	"os"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/libp2p/go-libp2p/core/peer"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/platform/web"
	"github.com/filecoin-project/lotus/api"
)

type LotusHealthAPI interface {
	NodeStatus(ctx context.Context, inclChainStatus bool) (api.NodeStatus, error)
	NetPeers(context.Context) ([]peer.AddrInfo, error)
	Version(context.Context) (api.APIVersion, error)
}

type Health struct {
	log   *logging.ZapEventLogger
	node  LotusHealthAPI
	build string
}

func NewHealth(log *logging.ZapEventLogger, node LotusHealthAPI, build string) *Health {
	return &Health{
		log:   log,
		node:  node,
		build: build,
	}
}

// Liveness returns status info if the service is alive.
func (h *Health) Liveness(w http.ResponseWriter, r *http.Request) {
	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	statusCode := http.StatusOK

	status, err := h.node.NodeStatus(r.Context(), true)
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	version, err := h.node.Version(r.Context())
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	h.log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	p, err := h.node.NetPeers(r.Context())
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	resp := data.LivenessResponse{
		Version:              version.String(),
		Epoch:                status.SyncStatus.Epoch,
		Behind:               status.SyncStatus.Behind,
		PeersToPublishMsgs:   status.PeerStatus.PeersToPublishMsgs,
		PeersToPublishBlocks: status.PeerStatus.PeersToPublishBlocks,
		PeerNumber:           len(p),
		Host:                 host,
		Build:                h.build,
	}

	if err := web.Respond(r.Context(), w, resp, http.StatusOK); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
}

// Readiness checks if the components are ready and if not will return a 500 status.
func (h *Health) Readiness(w http.ResponseWriter, r *http.Request) {
	ctx, cancel := context.WithTimeout(r.Context(), time.Second)
	defer cancel()

	status := "ok"
	statusCode := http.StatusOK

	if _, err := h.node.Version(ctx); err != nil {
		status = "lotus not ready"
		statusCode = http.StatusInternalServerError
	}

	h.log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	resp := struct {
		Status string `json:"status"`
	}{
		Status: status,
	}

	if err := web.Respond(r.Context(), w, resp, http.StatusOK); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
}
