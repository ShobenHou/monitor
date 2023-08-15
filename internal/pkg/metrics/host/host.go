package host

import (
	"time"
	"os"
	"github.com/ShobenHou/monitor/internal/pkg/metrics"
	"github.com/ShobenHou/monitor/internal/pkg/utils/json"
)

type HostStats struct {
}

func NewHostStats() *HostStats {
	return &HostStats{}
}

func (h *HostStats) Gather(acc *metrics.Accumulator) error {
	var err error

	measurement := "host"
	now := time.Now()

	var info *Info
	info, err = HostInfo()
	if err != nil {
		return err
	}
	acc.Add(measurement, map[string]string{
		"host": acc.Uid,
		"node": os.Getenv("MY_NODE_NAME"),
	}, json.StructToMap(info), now)

	return nil
}

func init() {
	metrics.Add("host", func() metrics.Metric {
		return NewHostStats()
	})
}
