package agent

import (
	"fmt"
	kafkaHelper "github.com/ShobenHou/monitor/internal/pkg/kafka"
	"github.com/ShobenHou/monitor/internal/pkg/metrics"
	_ "github.com/ShobenHou/monitor/internal/pkg/metrics/all"
	"time"
	//"github.com/ShobenHou/monitor/internal/pkg/rpc"
	log "github.com/sirupsen/logrus"
)

type Agent struct {
	Conf    *AgentConf
	restart chan bool //TODO:added
}

type AgentConf struct {
	Interval string   `yaml:"interval"`
	Metrics  []string `yaml:"metrics"`
	Addr     string   `yaml:"addr"`
}

func NewAgent(conf *AgentConf) *Agent {
	if conf != nil {
		return &Agent{
			Conf:    conf,
			restart: make(chan bool),
		}
	} else {
		return &Agent{
			Conf: &AgentConf{
				Interval: "10s",
				Metrics: []string{
					"cpu",
					"mem",
					"swap",
					"system",
					"kernel",
					"processes",
					"disk",
					"diskio",
					"load",
					"net",
				},
				Addr: "localhost:55555",
			},
			restart: make(chan bool),
		}
	}
}

func (a *Agent) UpdateConfig(newConf *kafkaHelper.MonitorConf) { // TODO:ADDED
	//newAgentConf := &AgentConf{Interval: newConf.MonitorInterval, Metrics: newConf.MonitorMetrics, Addr: a.Conf.Addr}
	a.Conf.Interval = newConf.MonitorInterval

	a.Conf.Metrics = newConf.MonitorMetrics
	//a.Conf = newAgentConf //TODO: map MonitorConfig values to AgentConf Values
	fmt.Println("Restarting")

	timeout := time.After(5 * time.Second) // 设置超时时间为3秒
	select {
	case a.restart <- true: // 尝试写入管道
		fmt.Println("Writing restart success")
	case <-timeout: // 超时处理
		fmt.Println("Writing resetart failed")
	}

}

func (a *Agent) Run(quit chan bool) error {
	var err error
	var interval time.Duration

	acc, err := metrics.NewAccumulator(100)
	if err != nil {
		return err
	}
	defer acc.Close() // will end report goruotines
	// Start report goruotines
	go acc.Report(3)

	interval, err = time.ParseDuration(a.Conf.Interval)
	if err != nil {
		return err
	}

	// Start collect goruotine
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-a.restart:
			log.Info("Restarting data collection and reporting goroutines")
			ticker.Stop()
			acc.Close()

			acc, err = metrics.NewAccumulator(100)
			if err != nil {
				return err
			}
			defer acc.Close()
			go acc.Report(3)

			interval, err = time.ParseDuration(a.Conf.Interval)
			if err != nil {
				return err
			}
			ticker = time.NewTicker(interval)

		case <-ticker.C:
			for _, m := range a.Conf.Metrics {
				creator, ok := metrics.MetricCreators[m]
				if !ok {
					return fmt.Errorf("metric %s is illegal", m)
				}
				metric := creator()
				if err := metric.Gather(acc); err != nil {
					log.Infof("gather error: %v", err)
				}
			}

		case <-quit:
			log.Info("agent collect loop quit")
			return nil
		}
	}
}
