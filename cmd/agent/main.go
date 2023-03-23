package main

import (
	"context"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"os"

	//    "gopkg.in/yaml.v2"
	log "github.com/sirupsen/logrus"

	"github.com/ShobenHou/monitor/internal/agent"
	"github.com/ShobenHou/monitor/internal/pkg/influxdb"
	"github.com/ShobenHou/monitor/internal/pkg/logger"
	//"github.com/githubzjm/tuo/internal/agent/plugin"
	// "github.com/githubzjm/tuo/internal/agent/config"
)

// cmdline flags
var (
	// -h flag is set by default, which outputs supported flags
	fVersion = flag.Bool("v", false, "display version and exit")
	//fConfigFile = flag.String("f", "./configs/agent.yml", "specify config file")
)

var (
	version string // set by ldflags in Makefile
)

func init() {
	if version == "" {
		version = "unknown"
	}
}

func main() {
	var err error
	var isUp bool

	// parse flag
	flag.Parse()
	switch {
	case *fVersion:
		fmt.Printf("agent v%v\n", version)
		return
	}

	// load config
	// var conf *config.Config
	// conf, err = config.LoadFile(*fConfigFile)
	// if err != nil {
	// 	fmt.Printf("load config file failed: %v\n", err)
	// 	os.Exit(1)
	// }
	// fmt.Println(conf)

	// init log
	logConf := &logger.LogConf{
		LogOutput: logger.LogOutputStdout,
	}
	logger.Setup(logConf)
	log.Info("logger init done")

	reqLogger := logger.GetReqLogger("reqid", "userip")
	reqLogger.Info("test msg")
	//log.Info("log to file")

	// test metric
	//influxdbWriteAPI := influxdbUtil.NewWriteAPI(influxdbClient, conf.Influxdb.Org, conf.Influxdb.Bucket)
	// collect
	// cpuStats := cpu.NewCPUStats()
	// times, _ := cpu.CollectCPUTimes()
	// m, _ := cpu.TimesStatToMap(times)

	// init influxdb
	influxdb.InitDB("http://127.0.0.1:8086", "init-token", "init-org", "init-bucket")
	defer influxdb.Close()
	if isUp, err = influxdb.Client.Ping(context.Background()); err != nil {
		log.Fatalf("ping influxdb failed: %v", err)
	} else {
		log.Infof("ping influxdb: %v", isUp)
	}

	// init agent
	agentConf := &agent.AgentConf{
		Interval: "1s",
		Metrics: []string{
			"cpu",
			"mem",
			"host",
		},
		Addr: "localhost:55555",
	}
	agent := agent.NewAgent(agentConf)

	//TODO：ADDED
	// Connect to Kafka
	kafkaBroker := "localhost:9092" // Replace with your Kafka broker(s) address
	kafkaGroupId := "agent-group"   // You may use a unique name for your agent group
	kafkaTopic := "agent-config"

	consumer, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers": kafkaBroker,
		"group.id":          kafkaGroupId,
		"auto.offset.reset": "earliest",
	})

	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
		os.Exit(1)
	}

	defer consumer.Close()

	// Subscribe to the configuration topic
	err = consumer.Subscribe(kafkaTopic, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to topic '%s': %v", kafkaTopic, err)
		os.Exit(1)
	}

	// Start listening for configuration updates
	go func() {
		for {
			msg, err := consumer.ReadMessage(-1)
			if err != nil {
				log.Errorf("Failed to read message from Kafka: %v", err)
				continue
			}

			newConf := &agent.AgentConf{}
			err = json.Unmarshal(msg.Value, newConf)
			if err != nil {
				log.Errorf("Failed to unmarshal configuration: %v", err)
				continue
			}

			agent.UpdateConfig(newConf)
			log.Infof("Updated agent configuration: %v", newConf)
		}
	}()
	//ADDED end

	quit := make(chan bool)
	defer func() {
		quit <- true // TODO need waitgroup to make sure the quit procedure is done,
		// otherwise main goroutine will end before the goroutines
	}()

	err = agent.Run(quit)
	if err != nil {
		log.Info(err)
	}

	// log.Info("sleep start")
	// time.Sleep(time.Second * 12)
	// log.Info("sleep end")

}
