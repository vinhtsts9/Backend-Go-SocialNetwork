package worker

import (
	"context"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/hibiken/asynq"
)

type TaskDistributor interface {
	DistributeTaskPost(
		ctx context.Context,
		payload *model.CreatePostInput,
		opts ...asynq.Option,
	) (database.GetPostByIdRow, error)
}

type RedisTaskDistributor struct {
	client *asynq.Client
	store  database.Store
}

func NewRedisTaskDistributor(redisOpt asynq.RedisClientOpt) TaskDistributor {
	client := asynq.NewClient(redisOpt)
	store := database.NewStore(global.MdbcHaproxy)
	return &RedisTaskDistributor{
		client: client,
		store:  store,
	}
}
