package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"

	"github.com/IBM/sarama"
)

type sPostProcessor struct {
	r *database.Queries
}

func NewPostProcessorImpl(r *database.Queries) *sPostProcessor {
	return &sPostProcessor{
		r: r,
	}
}

// ProcessPostMessage xử lý tin nhắn từ Kafka
func (s *sPostProcessor) ProcessPostMessage(ctx context.Context, topic string) {
	handler := func(message *sarama.ConsumerMessage) error {
		// Parse thông báo nhận được
		var post model.Post
		err := json.Unmarshal(message.Value, &post)
		if err != nil {
			global.Logger.Sugar().Errorf("Failed to unmarshal post: %v", err)
			return err
		}

		// Truy vấn người theo dõi
		followers, err := s.r.GetFollowersByUserId(ctx, sql.NullInt64{Int64: int64(post.UserID), Valid: true})
		if err != nil {
			global.Logger.Sugar().Errorf("Failed to get followers: %v", err)
			return err
		}

		onlineFollowers := []model.UserBase{}

		for _, follower := range followers {

			follower := model.UserBase{
				UserID:         follower.UserID,
				UserLogoutTime: follower.UserLogoutTime,
				UserState:      uint8(follower.CalculatedUserState),
			}
			if follower.UserState == 1 || follower.UserState == 2 {
				onlineFollowers = append(onlineFollowers, follower)
			} else if follower.UserState == 3 {
				query := database.UpdateUserStateParams{
					UserState: 3,
					UserID:    follower.UserID,
				}
				err := s.r.UpdateUserState(ctx, query)
				if err != nil {
					global.Logger.Sugar().Error("failed to update state for user %d, %v ", follower.UserID, err)
				}
			}
			if len(onlineFollowers) > 0 {
				if err := updateTimelineCache(post, onlineFollowers); err != nil {
					global.Logger.Sugar().Errorf("Failed to update timeline cache: %v", err)
					return err
				}
			}

		}
		return nil
	}
	err := global.KafkaConsumer.Consume(topic, handler)
	if err != nil {
		global.Logger.Sugar().Errorf("Failed to start consuming from topic %s: %v", topic, err)
	}
}

// Cập nhật cache timeline
func updateTimelineCache(post model.Post, onlineFollowers []model.UserBase) error {
	const batchsize = 100
	const maxWorkers = 5

	jobs := make(chan []model.UserBase, len(onlineFollowers)/batchsize+1)
	errChan := make(chan error, maxWorkers)
	worker := func() {
		for batch := range jobs {
			pipeline := global.Rdb.Pipeline()
			for _, follower := range batch {
				key := fmt.Sprintf("timeline%d", follower.UserID)

				timeline := getCache(key)
				timeline = append(timeline, post)

				data, err := json.Marshal(timeline)
				if err != nil {
					errChan <- fmt.Errorf("failed to marshal timeline %w", err)
					return
				}
				pipeline.Set(context.Background(), key, data, 0)
			}

			_, err := pipeline.Exec(context.Background())
			if err != nil {
				errChan <- fmt.Errorf("pipeline execution failed: %w", err)
				return
			}
		}
		errChan <- nil // bao thanh cong khi worker hoan thanh
	}
	for i := 0; i < len(onlineFollowers); i += batchsize {
		end := i + batchsize
		if end > len(onlineFollowers) {
			end = len(onlineFollowers)
		}
		batch := onlineFollowers[i:end]
		jobs <- batch
	}
	close(jobs)

	for i := 0; i < maxWorkers; i++ {
		go worker()
	}
	for i := 0; i < maxWorkers; i++ {
		if err := <-errChan; err != nil {
			return err
		}
	}
	return nil
}

// Lấy cache từ Redis
func getCache(key string) []model.Post {
	var timeline []model.Post
	data, err := global.Rdb.Get(context.Background(), key).Result()
	if err == nil {
		json.Unmarshal([]byte(data), &timeline)
	}
	return timeline
}

// Ghi cache vào Redis
func setCache(key string, timeline []model.Post) error {
	data, err := json.Marshal(timeline)
	if err != nil {
		return fmt.Errorf("failed to marshal timeline: %w", err)
	}
	return global.Rdb.Set(context.Background(), key, data, 0).Err()
}
