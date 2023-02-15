package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"log"
	"time"
)

type MonitoringConfig struct {
	MonitoringItems  []string      `json:"monitoring_items"`
	CollectionPeriod time.Duration `json:"collection_period"`
}

func main() {
	// Set up the Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer p.Close()

	topic := "monitoring_configurations"

	// Set up a ticker to generate monitoring configuration messages periodically
	ticker := time.NewTicker(5 * time.Second)
	defer ticker.Stop()

	// Loop through the ticker ticks and generate monitoring configuration messages
	for range ticker.C {
		config := MonitoringConfig{
			MonitoringItems:  []string{"cpu", "memory", "disk"},
			CollectionPeriod: 30 * time.Second,
		}

		// Convert the monitoring configuration to JSON
		configBytes, err := json.Marshal(config)
		if err != nil {
			log.Printf("Failed to marshal monitoring configuration: %v", err)
			continue
		}

		// Publish the monitoring configuration message to the Kafka topic
		err = p.Produce(&kafka.Message{
			TopicPartition: kafka.TopicPartition{Topic: &topic, Partition: kafka.PartitionAny},
			Value:          configBytes,
		}, nil)
		if err != nil {
			log.Printf("Failed to publish monitoring configuration message: %v", err)
			continue
		}
		fmt.Print(fmt.Printf("Sent monitoring configuration: %+v\n", config))
	}

}
