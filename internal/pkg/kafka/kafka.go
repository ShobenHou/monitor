package kafka

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
)

type MonitorConf struct {
	MonitorMetrics  []string `json:"monitor_metrics"`
	MonitorInterval string   `json:"monitor_interval"`
}

var kafkaProducer *kafka.Producer

func init() {
	// Initialize the Kafka producer.
	var err error
	kafkaProducer, err = kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "kafka:9092"})
	if err != nil {
		panic(err)
	}
}

func PublishConfigToKafka(monitorConf MonitorConf) error {
	monitorConfJSON, err := json.Marshal(monitorConf)
	if err != nil {
		return err
	}

	kafkaTopic := "monitoring_configurations"
	msg := &kafka.Message{
		TopicPartition: kafka.TopicPartition{Topic: &kafkaTopic, Partition: kafka.PartitionAny},
		//TODO: Key:            []byte(agentID),
		Value: monitorConfJSON,
	}

	deliveryChan := make(chan kafka.Event)
	defer close(deliveryChan)

	fmt.Printf("PublishConfigToKafka--Producing MonitorConfig to topic:%s\n", msg.TopicPartition.Topic)
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
