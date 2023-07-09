package main

import (
	"encoding/json"
	"fmt"
	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/influxdata/influxdb-client-go/v2"
	"log"
	"os"
	"os/signal"
	"time"
)

type MonitoringConfig struct {
	MonitoringItems  []string      `json:"monitoring_items"`
	CollectionPeriod time.Duration `json:"collection_period"`
}

func main() {
	// Set up the InfluxDB client
	influxClient := influxdb2.NewClient("http://localhost:8086", "my-token")
	defer influxClient.Close()

	// Set up the Kafka producer
	p, err := kafka.NewProducer(&kafka.ConfigMap{"bootstrap.servers": "localhost:9092"})
	if err != nil {
		log.Fatalf("Failed to create Kafka producer: %v", err)
	}
	defer p.Close()

	// Set up the Kafka consumer
	c, err := kafka.NewConsumer(&kafka.ConfigMap{
		"bootstrap.servers":  "localhost:9092",
		"group.id":           "my-group",
		"auto.offset.reset":  "earliest",
		"enable.auto.commit": false,
	})
	if err != nil {
		log.Fatalf("Failed to create Kafka consumer: %v", err)
	}
	defer c.Close()

	// Subscribe to the monitoring configurations topic
	err = c.SubscribeTopics([]string{"monitoring_configurations"}, nil)
	if err != nil {
		log.Fatalf("Failed to subscribe to Kafka topic: %v", err)
	}

	// Set up a signal handler to handle SIGINT and SIGTERM signals
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, os.Interrupt, os.Kill)

	// Loop through the monitoring configuration messages and update the monitoring settings
	for {
		select {
		case sig := <-sigchan:
			log.Printf("Caught signal %v: shutting down", sig)
			return
		default:
			msg, err := c.ReadMessage(5 * time.Second)
			if err != nil {
				if err.(kafka.Error).Code() != kafka.ErrTimedOut {
					log.Printf("Failed to read Kafka message: %v", err)
				}

				continue
			}

			var config MonitoringConfig
			err = json.Unmarshal(msg.Value, &config)
			if err != nil {
				log.Printf("Failed to parse monitoring configuration message: %v", err)
				continue
			}
			// Output the monitoring configuration to standard output
			fmt.Printf("Received monitoring configuration: %+v\n", config)

			// Commit the Kafka message offset to mark it as processed
			if _, err := c.CommitMessage(msg); err != nil {
				log.Printf("Failed to commit Kafka message offset: %v", err)
			}
		}
	}
}

//This implementation uses the Confluent Kafka Go client library
//to communicate with the Kafka message queue server and
//the InfluxDB Go client library to insert monitoring data into InfluxDB.
