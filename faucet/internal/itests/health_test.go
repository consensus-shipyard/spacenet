package itests

import (
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"

	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/faucet/internal/data"
	handler "github.com/filecoin-project/faucet/internal/http"
	"github.com/filecoin-project/faucet/internal/itests/kit"
)

type HealthTests struct {
	handler http.Handler
}

func Test_Health(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotus()
	srv := handler.HealthHandler(log, lotus, "build")
	tests := HealthTests{
		handler: srv,
	}
	t.Run("liveness", tests.liveness)
	t.Run("readiness", tests.readiness)
}

func (ht *HealthTests) liveness(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	var resp data.LivenessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, uint64(10), resp.Epoch)
	require.Equal(t, 1, resp.PeerNumber)
}

func (ht *HealthTests) readiness(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}
