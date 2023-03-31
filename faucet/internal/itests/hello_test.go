package itests

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"testing"

	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/faucet/internal/data"
	handler "github.com/filecoin-project/faucet/internal/http"
	"github.com/filecoin-project/faucet/internal/itests/kit"
)

type HelloTests struct {
	handler http.Handler
}

func Test_Hello(t *testing.T) {
	log := logging.Logger("TEST-HELLO-SERVICE")

	lotus := kit.NewFakeLotus()

	srv := handler.HelloHandler(log, lotus)

	tests := HelloTests{
		handler: srv,
	}

	t.Run("hello", tests.hello)
}

func (ht *HelloTests) hello(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/hello", nil)
	w := httptest.NewRecorder()

	ht.handler.ServeHTTP(w, r)

	require.Equal(t, http.StatusOK, w.Code)
	fmt.Println(w.Body)

	var resp data.HelloResponse
	err := json.Unmarshal(w.Body.Bytes(), &resp)
	require.NoError(t, err)
	require.Equal(t, uint64(10), resp.Epoch)
	require.Equal(t, 1, resp.PeerNumber)
}
