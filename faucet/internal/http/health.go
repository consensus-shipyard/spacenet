package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
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
	ID(context.Context) (peer.ID, error)
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
	id, err := h.node.ID(r.Context())
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
		PeerID:               id.String(),
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

	daemonStatus := "ok"
	validatorStatus := "ok"
	statusCode := http.StatusOK

	if _, err := h.node.Version(ctx); err != nil {
		daemonStatus = "lotus not ready"
		statusCode = http.StatusInternalServerError
	}

	if err := h.checkValidatorStatus(); err != nil {
		validatorStatus = "validator not ready"
		statusCode = http.StatusInternalServerError
	}

	h.log.Infow("readiness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)

	resp := struct {
		DaemonStatus    string `json:"daemon_status"`
		ValidatorStatus string `json:"validator_status"`
	}{
		DaemonStatus:    daemonStatus,
		ValidatorStatus: validatorStatus,
	}

	if err := web.Respond(r.Context(), w, resp, http.StatusOK); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Health) checkValidatorStatus() error {
	grep := exec.Command("grep", "[e]udico mir validator")
	ps := exec.Command("ps", "ax")

	pipe, _ := ps.StdoutPipe()
	defer func(pipe io.ReadCloser) {
		err := pipe.Close()
		if err != nil {
			h.log.Infow("checkValidatorStatus error", err)
		}
	}(pipe)

	grep.Stdin = pipe
	if err := ps.Start(); err != nil {
		return err
	}

	// Run and get the output of grep.
	o, err := grep.Output()
	if err != nil {
		return err
	}
	if o == nil {
		return fmt.Errorf("validator not found")
	}
	return nil
}
