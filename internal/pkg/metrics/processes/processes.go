package processes

import (
	"context"
	"time"
	"os"
	"fmt"
	"github.com/ShobenHou/monitor/internal/pkg/metrics"

	"github.com/shirou/gopsutil/v3/process"
)

type ProcStats struct{}

func NewProcStats() *ProcStats {
	return &ProcStats{}
}

func (c *ProcStats) Gather(acc *metrics.Accumulator) error {
	// Get all running processes
	processes, err := process.ProcessesWithContext(context.Background())
	if err != nil {
		return err
	}

	now := time.Now()
	for _, p := range processes {
		if p.Pid == 0 {
			continue
		}
		fields := make(map[string]interface{})
		tags := map[string]string{
			"pid":  fmt.Sprintf("%d", p.Pid),
			"node": os.Getenv("MY_NODE_NAME"),
		}
		// CPU Percent
		cpuPercent, err := p.CPUPercentWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-CPUPercentWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["cpu_percent"] = cpuPercent
		}
		// Memory Percent
		memPercent, err := p.MemoryPercentWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-MemoryPercentWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["memory_percent"] = memPercent
		}
		// Process Status
		status, err := p.StatusWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-StatusWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["status"] = status
		}
		// Number of Threads
		numThreads, err := p.NumThreadsWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-NumThreadsWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["num_threads"] = numThreads
		}
		// Number of File Descriptors
		numFDs, err := p.NumFDsWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-NumFDsWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["num_fds"] = numFDs
		}
		// Number of Network Connections
		connections, err := p.ConnectionsMaxWithContext(context.Background(), 0)
		if err != nil {
			fmt.Println("ERROR-ConnectionsMaxWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["num_connections"] = len(connections)
		}
		// Process Create Time
		createTime, err := p.CreateTimeWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-CreateTimeWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["create_time"] = createTime
		}
		// Process Command Line
		cmdline, err := p.CmdlineWithContext(context.Background())
		if err != nil {
			fmt.Println("ERROR-CmdlineWithContext")
			fmt.Println("error gathering metrics for PID", p.Pid, ":", err)
		} else {
			fields["cmdline"] = cmdline
		}
		// Add the metrics to the accumulator
		acc.Add("processes", tags, fields, now)
	}

	return nil
}

func init() {
	metrics.Add("processes", func() metrics.Metric {
		return NewProcStats()
	})
}
