package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"time"

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

// GetAllPosts lấy danh sách tất cả bài viết
func (s *sPostProcessor) GetAllPosts(ctx context.Context) (codeRs int, data []model.Post, err error) {
	rows, err := s.r.GetAllpost(ctx)
	if err != nil {
		return response.ErrCodePostFailed, nil, err
	}
	var posts []model.Post
	for _, row := range rows {
		post := model.Post{
			ID:        uint32(row.ID),
			UserID:    uint32(row.UserID),
			Title:     row.Title,
			Content:   string(row.Content),
			CreatedAt: row.CreatedAt.Time.String(),
			UpdatedAt: row.UpdatedAt.Time.String(),
		}
		posts = append(posts, post)
	}
	return response.ErrCodeSuccess, posts, nil
}

// GetPostById lấy bài viết theo ID
func (s *sPostProcessor) GetPostById(ctx context.Context, postId string) (codeRs int, data model.Post, err error) {
	id, err := parsePostId(postId)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	row, err := s.r.GetPostById(ctx, id)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	post := model.Post{
		ID:        uint32(row.ID),
		UserID:    uint32(row.UserID),
		Title:     row.Title,
		Content:   string(row.Content),
		CreatedAt: row.CreatedAt.Time.String(),
		UpdatedAt: row.UpdatedAt.Time.String(),
	}
	return response.ErrCodeSuccess, post, nil
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

		currentTime := time.Now()
		onlineFollowers := []model.UserBase{}

		for _, follower := range followers {
			var newState int
			follower := model.UserBase{
				UserID:         follower.UserID,
				UserLogoutTime: follower.UserLogoutTime,
				UserState:      follower.UserState,
			}
			if follower.UserState == 2 {
				if follower.UserLogoutTime.Valid {
					duration := currentTime.Sub(follower.UserLogoutTime.Time)
					if duration.Minutes() >= 10 {
						newState = 3
					} else {
						newState = 2
					}
				} else {
					newState = 3
				}
			} else if follower.UserState == 1 {
				newState = 1
			}
			if newState == 1 || newState == 2 {
				onlineFollowers = append(onlineFollowers, follower)
			} else if newState == 3 {
				query := database.UpdateUserStateParams{
					UserState: 3,
					UserID:    follower.UserID,
				}
				err := s.r.UpdateUserState(ctx, query)
				if err != nil {
					global.Logger.Sugar().Errorf("Failed to update user state for user %d: %v", follower.UserID, err)
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
	for _, follower := range onlineFollowers {
		key := fmt.Sprintf("timeline:%d", follower.UserID)
		timeline := getCache(key)
		timeline = append(timeline, post)

		// Lưu lại vào Redis
		err := setCache(key, timeline)
		if err != nil {
			return fmt.Errorf("failed to set cache for user %d: %w", follower.UserID, err)
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
