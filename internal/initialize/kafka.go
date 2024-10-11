package initialize

import (
	"go-ecommerce-backend-api/m/v2/global"
	"log"

	kafka "github.com/segmentio/kafka-go"
)

var KafkaProducer *kafka.Writer

func InitKafKa() {
	global.KafkaProducer = &kafka.Writer{
		Addr:     kafka.TCP("127.0.0.1:9092"),
		Topic:    "otp-auth-topic",
		Balancer: &kafka.LeastBytes{},
	}
}

func CloseKafka() {
	if err := global.KafkaProducer.Close(); err != nil {
		log.Fatalf("Failed to close kafka producer %v", err)
	}
}
