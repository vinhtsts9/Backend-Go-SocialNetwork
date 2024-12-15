// initialize/kafka.go
package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/package/kafka"
	"log"
)

func InitKafka() {
	brokers := []string{"localhost:9092"}
	producer, err := kafka.GetProducer(brokers)
	if err != nil {
		log.Fatalf("Error initializing Kafka producer: %v", err)
	}
	global.KafkaProducer = producer

	global.KafkaConsumer, err = kafka.NewConsumer(brokers)
	if err != nil {
		log.Fatalf("Error initializing Kafka consumer: %v", err)
	}
}
