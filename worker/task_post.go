package worker

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/internal/service/impl"

	"github.com/hibiken/asynq"

	"github.com/rs/zerolog/log"
)

func (d *RedisTaskDistributor) DistributeTaskPost(
	ctx context.Context,
	payload *model.CreatePostInput,
	opts ...asynq.Option,
) (database.GetPostByIdRow, error) {
	global.Logger.Sugar().Infof("Model create Post: %+v", payload)

	// Serialize image_paths để lưu vào database
	contentJSON, err := json.Marshal(payload.ImagePaths)
	if err != nil {
		return database.GetPostByIdRow{}, fmt.Errorf("failed to marshal image paths: %w", err)
	}

	// Tạo đối tượng CreatePostParams để lưu vào DB
	createPostParams := database.CreatePostParams{
		Title:        payload.Title,
		ImagePaths:   contentJSON,
		UserNickname: payload.UserNickname,
		UserID:       payload.UserId,
	}

	// Lưu vào database
	result, err := d.store.CreatePost(ctx, createPostParams)
	if err != nil {
		return database.GetPostByIdRow{}, fmt.Errorf("failed to create post: %w", err)
	}
	lastId, err := result.LastInsertId()

	// Lấy thông tin bài post theo id vừa lấy được
	dataRs, err := d.store.GetPostById(ctx, uint64(lastId))
	if err != nil {
		return database.GetPostByIdRow{}, fmt.Errorf("failed to get post by id: %w", err)
	}

	// Marshal payload cho task asynq
	jsonPayload, err := json.Marshal(dataRs)
	if err != nil {
		return database.GetPostByIdRow{}, fmt.Errorf("failed to marshal task payload: %w", err)
	}

	task := asynq.NewTask(TaskSendPost, jsonPayload, opts...)

	info, err := d.client.EnqueueContext(ctx, task)
	if err != nil {
		return database.GetPostByIdRow{}, fmt.Errorf("failed to enqueue task: %w", err)
	}

	log.Info().
		Str("type", task.Type()).
		Bytes("payload", task.Payload()).
		Str("queue", info.Queue).
		Int("max_retry", info.MaxRetry).
		Msg("enqueued task")

	return dataRs, nil
}

func (p *RedisTaskProcessor) ProcessTaskSendPost(ctx context.Context, task *asynq.Task) error {
	var payload model.Post
	if err := json.Unmarshal(task.Payload(), &payload); err != nil {
		return fmt.Errorf("failed to unmarshal payload: %w", asynq.SkipRetry)
	}

	// Truy vấn người theo dõi
	followers, err := p.store.GetFollowersByUserId(ctx, sql.NullInt64{Int64: int64(payload.UserId), Valid: true})
	if err != nil {
		global.Logger.Sugar().Errorf("Failed to get followers: %v", err)
		return err
	}
	onlineFollowers := []model.UserBase{}
	for _, follower := range followers {
		userBase := model.UserBase{
			UserID:         follower.UserID,
			UserLogoutTime: follower.UserLogoutTime,
			UserState:      uint8(follower.CalculatedUserState),
		}

		// Người dùng online
		if userBase.UserState == 1 || userBase.UserState == 2 {
			onlineFollowers = append(onlineFollowers, userBase)
		}
		// Người dùng offline lâu
		if userBase.UserState == 3 {
			query := database.UpdateUserStateParams{
				UserState: 3,
				UserID:    userBase.UserID,
			}
			if err := p.store.UpdateUserState(ctx, query); err != nil {
				global.Logger.Sugar().Errorf("Failed to update state for user %d, %v", userBase.UserID, err)
			}
		}
	}

	// Cập nhật cache cho những người theo dõi online
	if len(onlineFollowers) > 0 {
		if err := impl.UpdateTimelineCache(payload, onlineFollowers); err != nil {
			global.Logger.Sugar().Errorf("Failed to update timeline cache: %v", err)
			return err
		}
	}

	global.Logger.Sugar().Infof("Processed post task for user %d", payload.UserId)
	return nil
}
