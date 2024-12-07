package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

var KafkaProducer *kafka.Writer
var KafkaConsumer *kafka.Reader

func InitKafKa() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("localhost:9092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
	global.KafkaConsumer = kafka.NewReader(kafka.ReaderConfig{
		Brokers: []string{"localhost:9092"},
		Topic:   "room",
		GroupID: "chat-consumer-group",
	})
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer: %v", err)
	}
}
