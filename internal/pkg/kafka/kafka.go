package kafka

import (
	"encoding/json"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"time"
)

type MonitorConf struct {
	// TODO: different from AgentConf in internal/agent/agent.go
	MonitorMetrics  []string      `json:"monitor_metrics"`
	MonitorInterval time.Duration `json:"monitor_interval"`
}

var kafkaProducer *kafka.Producer

func init() {
	// Initialize the Kafka producer.
	var err error
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		panic(err)
	}
}

func PublishConfigToKafka(agentID string, monitorConf MonitorConf) error {
	monitorConfJSON, err := json.Marshal(monitorConf)
	if err != nil {
		return err
	}

	kafkaTopic := "agent-config"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		Key:            []byte(agentID),
		Value:          monitorConfJSON,
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	err = kafkaProducer.Produce(msg, deliveryChan)
	if err != nil {
		return err
	}

	e := <-deliveryChan
	m := e.(*kafka.Message)

	if m.TopicPartition.Error != nil {
		return m.TopicPartition.Error
	}

	return nil
}
