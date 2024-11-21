package sendto

import (
	"context"
	"encoding/json"
	"go-ecommerce-backend-api/m/v2/global"
	"time"

	"github.com/segmentio/kafka-go"
)

func SendEmailKafka(email []string, otp string) error {
	body := make(map[string]interface{})

	body["otp"] = otp
	body["email"] = email

	bodyRequest, _ := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}

	err := global.KafkaProducer.WriteMessages(context.Background(), message)
	if err != nil {
		return err
	}
	return nil
}
