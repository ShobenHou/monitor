package agent

import (
	"fmt"
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
			Conf: conf,
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
				},
				Addr: "localhost:55555",
			},
		}
	}
}

func (a *Agent) UpdateConfig(newConf *AgentConf) { // TODO:ADDED
	a.Conf = newConf
	a.restart <- true
}

/*
func (a *Agent) subscribeKafka(kafkaAddr, topic string, updateConfig func(*AgentConf)) {//TODO：ADDED
	r := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaAddr},
		Topic:     topic,
		Partition: 0,
		MinBytes:  10e3,
		MaxBytes:  10e6,
	})

	for {
		m, err := r.ReadMessage(context.Background())
		if err != nil {
			log.Errorf("Error reading message from Kafka: %v", err)
			continue
		}

		var newConf AgentConf
		if err := json.Unmarshal(m.Value, &newConf); err != nil {
			log.Errorf("Error unmarshalling Kafka message: %v", err)
			continue
		}

		updateConfig(&newConf)
		log.Infof("Updated configuration: %v", newConf)
	}
}
*/

func (a *Agent) Run(quit chan bool) error {
	var err error
	var interval time.Duration

	// Add the following line to start the Kafka subscription goroutine
	//go a.subscribeKafka("localhost:9092", "agent-config", a.UpdateConfig) //TODO: ADDED//放的位置不一定对

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
		case <-a.restart: //TODO: added
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

		case <-quit:
			log.Info("agent collect loop quit")
			return nil
		}
	}
}
