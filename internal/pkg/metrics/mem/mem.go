package mem

import (
	"time"
	"os"
	"github.com/ShobenHou/monitor/internal/pkg/metrics"
	"github.com/ShobenHou/monitor/internal/pkg/utils/json"
)

type MemStats struct {
}

func NewMemStats() *MemStats {
	return &MemStats{}
}

func (m *MemStats) Gather(acc *metrics.Accumulator) error {
	var err error

	measurement := "mem"
	now := time.Now()

	// Swap mem
	var swap *Swap
	swap, err = SwapMem()
	if err != nil {
		return err
	}
	acc.Add(measurement, map[string]string{
		"node": os.Getenv("MY_NODE_NAME"),
	}, json.StructToMap(swap), now)

	// Virtual mem
	var virtual *Virtual
	virtual, err = VirtualMem()
	if err != nil {
		return err
	}
	acc.Add(measurement, map[string]string{
		"node": os.Getenv("MY_NODE_NAME"),
	}, json.StructToMap(virtual), now)

	return nil
}

func init() {
	metrics.Add("mem", func() metrics.Metric {
		return NewMemStats()
	})
}
