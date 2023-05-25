package itests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/faucet/internal/data"
	"github.com/filecoin-project/faucet/internal/failure"
	handler "github.com/filecoin-project/faucet/internal/http"
	"github.com/filecoin-project/faucet/internal/itests/kit"
)

type HealthTests struct {
	handler http.Handler
}

func TestValidatorHealth(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotusNoCrash()
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)
	srv := handler.HealthHandler(log, lotus, d, "build")
	test := HealthTests{
		handler: srv,
	}
	t.Run("validator-liveness", test.livenessForValidator)
	t.Run("validator-readiness", test.readinessForValidator)
}

func (ht *HealthTests) livenessForValidator(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	var resp data.LivenessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, 1, resp.PeerNumber)
}

func (ht *HealthTests) readinessForValidator(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestBootstrapHealth(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotusNoCrash()
	check := func() error {
		return fmt.Errorf("failed")
	}
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)
	srv := handler.HealthHandler(log, lotus, d, "build", check)
	test := HealthTests{
		handler: srv,
	}
	t.Run("bootstrap-liveness", test.livenessForBootstrap)
	t.Run("bootstrap-readiness", test.readinessForBootstrap)
}

func (ht *HealthTests) livenessForBootstrap(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/liveness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
	var resp data.LivenessResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, 1, resp.PeerNumber)
}

func (ht *HealthTests) readinessForBootstrap(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/readiness?bootstrap=true", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusOK, w.Code)
}

func TestValidatorFailedHealth(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotusNoCrash()
	check := func() error {
		return fmt.Errorf("failed")
	}
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)
	srv := handler.HealthHandler(log, lotus, d, "build", check)
	test := HealthTests{
		handler: srv,
	}
	t.Run("failed-validator-readiness", test.failedReadinessForBootstrap)
}

func (ht *HealthTests) failedReadinessForBootstrap(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}

func TestValidatorFailedHealthWithFailedLotus(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotusWithFailedVersion()
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)
	srv := handler.HealthHandler(log, lotus, d, "build")
	test := HealthTests{
		handler: srv,
	}
	t.Run("failed-validator-readiness-failed-lotus", test.failedReadinessForBootstrapWithFailedLotus)
}

func (ht *HealthTests) failedReadinessForBootstrapWithFailedLotus(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/readiness", nil)
	w := httptest.NewRecorder()
	ht.handler.ServeHTTP(w, r)
	require.Equal(t, http.StatusInternalServerError, w.Code)
}
