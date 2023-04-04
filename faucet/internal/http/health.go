package http

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"

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
	check ValidatorHealthCheck
}

type ValidatorHealthCheck func() error

func NewHealth(log *logging.ZapEventLogger, node LotusHealthAPI, build string, check ...ValidatorHealthCheck) *Health {
	h := Health{
		log:   log,
		node:  node,
		build: build,
	}
	if check == nil {
		h.check = defaultValidatorHealthCheck
	} else {
		h.check = check[0]
	}
	return &h
}

// Liveness returns status info if the service is alive.
func (h *Health) Liveness(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()

	host, err := os.Hostname()
	if err != nil {
		host = "unavailable"
	}

	statusCode := http.StatusOK

	status, err := h.node.NodeStatus(ctx, true)
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	version, err := h.node.Version(ctx)
	if err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}

	h.log.Infow("liveness", "statusCode", statusCode, "method", r.Method, "path", r.URL.Path, "remoteaddr", r.RemoteAddr)

	p, err := h.node.NetPeers(ctx)
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
	ctx := r.Context()

	h.log.Infow("readiness", "method", r.Method, "path", r.URL.Path, "remote", r.RemoteAddr)

	ready := true
	if _, err := h.node.Version(ctx); err != nil {
		h.log.Infow("failed to connect to daemon", "readiness", "error", err)
		ready = false
	}

	// A node can be a bootstrap node or validator node. Bootstrap nodes run daemons only.
	// We signal that a node is a bootstrap node by accessing /readiness endpoint with "boostrap" parameter.
	isBootstrap := r.URL.Query().Get("bootstrap") != ""
	fmt.Println(isBootstrap)

	if !isBootstrap {
		if err := h.checkValidatorStatus(); err != nil {
			h.log.Infow("failed to connect to validator", "readiness", "error", err)
			ready = false
		}
	}

	if !ready {
		http.Error(w, http.StatusText(http.StatusInternalServerError), http.StatusInternalServerError)
		return
	}

	resp := struct {
		Status string `json:"status"`
	}{
		Status: "ok",
	}

	if err := web.Respond(ctx, w, resp, http.StatusOK); err != nil {
		web.RespondError(w, http.StatusInternalServerError, err)
		return
	}
}

func (h *Health) checkValidatorStatus() error {
	return h.check()
}

func defaultValidatorHealthCheck() error {
	grep := exec.Command("grep", "[e]udico mir validator")
	ps := exec.Command("ps", "ax")

	pipe, _ := ps.StdoutPipe()
	defer func(pipe io.ReadCloser) {
		pipe.Close() // nolint
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
