package initialize

import (
	"context"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"

	"github.com/hibiken/asynq"
	"github.com/redis/go-redis/v9"
	"go.uber.org/zap"
)

var ctx = context.Background()

func InitRedis() {
	r := global.Config.Redis
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:%v", r.Host, r.Port),
		Password: r.Password,
		DB:       r.Database,
		PoolSize: 10,
	})
	global.RedisOpt = asynq.RedisClientOpt{
		Addr: fmt.Sprintf("%s:%v", r.Host, r.Port),
	}
	_, err := rdb.Ping(ctx).Result()
	if err != nil {
		global.Logger.Error("Redis initialization Error", zap.Error(err))
	}

	global.Rdb = rdb
}
