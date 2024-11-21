package global

import (
	"database/sql"
	"go-ecommerce-backend-api/m/v2/package/setting"

	"github.com/redis/go-redis/v9"
	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config        setting.Config
	Logger        *zap.Logger
	Mdb           *gorm.DB
	Rdb           *redis.Client
	Mdbc          *sql.DB
	KafkaProducer *kafka.Writer
)
