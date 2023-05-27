package failure

import (
	"context"
	"fmt"
	"sync"
	"time"

	logging "github.com/ipfs/go-log/v2"

	"github.com/filecoin-project/faucet/internal/platform/lotus"
)

type Detector struct {
	m                         sync.Mutex
	log                       *logging.ZapEventLogger
	lastBlockHeight           uint64
	lotus                     lotus.API
	lastBlockHeightUpdateTime time.Time
	threshold                 time.Duration
	checkInterval             time.Duration
	stopChan                  chan bool
	ticker                    *time.Ticker
}

// NewDetector creates a new failure detector that checks height value each checkInterval
// and triggers a failure if there is no block height update in threshold.
func NewDetector(log *logging.ZapEventLogger, api lotus.API, checkInterval, threshold time.Duration) *Detector {
	d := Detector{
		checkInterval:             checkInterval,
		lotus:                     api,
		log:                       log,
		threshold:                 threshold,
		stopChan:                  make(chan bool),
		lastBlockHeightUpdateTime: time.Now(),
		ticker:                    time.NewTicker(checkInterval),
	}

	go d.run()

	return &d
}

func (d *Detector) run() {
	ctx, cancel := context.WithCancel(context.Background())

	for {
		select {
		case <-d.stopChan:
			cancel()
			close(d.stopChan)
			d.log.Infow("shutdown", "status", "detector stopped")
			return
		case <-d.ticker.C:
			status, err := d.lotus.NodeStatus(ctx, true)
			if err != nil {
				d.log.Errorw("error", "detector", "unable to get block", err)
			} else {
				height := status.SyncStatus.Epoch
				d.m.Lock()
				if d.lastBlockHeight != height {
					d.lastBlockHeight = height
					d.lastBlockHeightUpdateTime = time.Now()
				}
				d.m.Unlock()
			}
		}
	}

}

func (d *Detector) Stop() {
	d.ticker.Stop()
	close(d.stopChan)
}

func (d *Detector) GetLastBlockHeight() uint64 {
	d.m.Lock()
	defer d.m.Unlock()
	return d.lastBlockHeight
}

func (d *Detector) CheckProgress() error {
	d.m.Lock()
	defer d.m.Unlock()

	if time.Since(d.lastBlockHeightUpdateTime) > d.threshold {
		return fmt.Errorf("no blocks since block %d at %s",
			d.lastBlockHeight, d.lastBlockHeightUpdateTime.Format("2006-01-02 15:04:05"))
	}
	return nil
}
