package all

import (
	// Blank imports for plugins to register themselves
	// init func will be run automatically when package is imported
	_ "github.com/ShobenHou/monitor/internal/pkg/metrics/cpu"
	_ "github.com/ShobenHou/monitor/internal/pkg/metrics/heartbeat"
	_ "github.com/ShobenHou/monitor/internal/pkg/metrics/host"
	_ "github.com/ShobenHou/monitor/internal/pkg/metrics/mem"
)
