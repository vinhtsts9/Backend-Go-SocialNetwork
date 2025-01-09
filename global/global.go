package global

import (
	"database/sql"
	"go-ecommerce-backend-api/m/v2/package/cloudinary"
	"go-ecommerce-backend-api/m/v2/package/kafka"
	"go-ecommerce-backend-api/m/v2/package/setting"

	"github.com/casbin/casbin/v2"
	"github.com/elastic/go-elasticsearch/v8"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
	"gorm.io/gorm"
)

var (
	Config        setting.Config
	Logger        *zap.Logger
	Mdb           *gorm.DB
	Rdb           *redis.Client
	MdbcHaproxy   *sql.DB
	MdbcSlave     *sql.DB
	MdbcSlave2    *sql.DB // Added for MySQL Slave2
	MdbcSlave3    *sql.DB // Added for MySQL Slave3
	KafkaProducer *kafka.Producer
	KafkaConsumer *kafka.Consumer
	Casbin        *casbin.Enforcer
	Cloudinary    *cloudinary.CloudinaryService
	Elastic       *elasticsearch.Client
)
