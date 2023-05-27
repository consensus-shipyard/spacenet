package itests

import (
	"testing"
	"time"

	logging "github.com/ipfs/go-log/v2"
	"github.com/stretchr/testify/require"

	"github.com/filecoin-project/faucet/internal/failure"
	"github.com/filecoin-project/faucet/internal/itests/kit"
)

func TestDetectorWhenBlockProduced(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotusNoCrash()
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)

	lastBlockHeight := d.GetLastBlockHeight()

	for i := 0; i < 20; i++ {
		time.Sleep(300 * time.Millisecond)
		require.NoError(t, d.CheckProgress())
		h := d.GetLastBlockHeight()
		require.Greater(t, h, lastBlockHeight)
		lastBlockHeight = h
	}
}

func TestDetectorWhenCrash(t *testing.T) {
	log := logging.Logger("TEST-HEALTH")
	lotus := kit.NewFakeLotus(true, 10)
	d := failure.NewDetector(log, lotus, 100*time.Millisecond, time.Second)

	lastBlockHeight := d.GetLastBlockHeight()

	for i := 0; i < 20; i++ {
		time.Sleep(300 * time.Millisecond)
		h := d.GetLastBlockHeight()
		if h < 10 {
			require.NoError(t, d.CheckProgress())
			require.Greater(t, h, lastBlockHeight)
		} else if h > 10 {
			require.Error(t, d.CheckProgress())
			require.Equal(t, uint64(10), h)
		}
		lastBlockHeight = h
	}
}
