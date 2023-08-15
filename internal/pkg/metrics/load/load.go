package load

import (
	"context"
	"fmt"
	"os"
	"time"

	"github.com/ShobenHou/monitor/internal/pkg/metrics"

	"github.com/shirou/gopsutil/v3/load"
)

type LoadStats struct{}

func NewLoadStats() *LoadStats {
	return &LoadStats{}
}

func (c *LoadStats) Gather(acc *metrics.Accumulator) error {
	avgStat, err := load.AvgWithContext(context.Background())
	tags := map[string]string{
		"node": os.Getenv("MY_NODE_NAME"),
	}
	if err != nil {
		fmt.Println("ERROR-load.AvgWithContext")
		fmt.Printf("error getting load average: %w\n", err)
	} else {
		fields := map[string]interface{}{
			"load1":  avgStat.Load1,
			"load5":  avgStat.Load5,
			"load15": avgStat.Load15,
		}
		acc.Add("load", tags, fields, time.Now())
	}

	miscStat, err := load.MiscWithContext(context.Background())
	if err != nil {
		fmt.Println("ERROR-load.MiscWithContext")
		fmt.Printf("error getting miscellaneous statistics: %w\n", err)
	} else {
		fields := map[string]interface{}{
			"procs_total":   miscStat.ProcsTotal,
			"procs_created": miscStat.ProcsCreated,
			"procs_running": miscStat.ProcsRunning,
			"procs_blocked": miscStat.ProcsBlocked,
			"ctxt":          miscStat.Ctxt,
		}
		acc.Add("load_misc", tags, fields, time.Now())
	}

	return nil
}

func init() {
	metrics.Add("load", func() metrics.Metric {
		return NewLoadStats()
	})
}
