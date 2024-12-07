package impl

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
)

type sPostProcessor struct {
	r *database.Queries
}

func NewPostProcessorImpl(r *database.Queries) *sPostProcessor {
	return &sPostProcessor{
		r: r,
	}
}

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

func (s *sPostProcessor) ProcessPostMessage(ctx context.Context) {
	for {
		msg, err := global.KafkaConsumer.ReadMessage(ctx)
		if err != nil {
			global.Logger.Sugar().Error("Failed to read message: %v", err)
			continue
		}

		var post model.Post
		err = json.Unmarshal(msg.Value, &post)
		if err != nil {
			global.Logger.Sugar().Error("Failed to Unmarshal post: %v", err)
		}

		// queries nguoi theo doi
		//followers := s.r.queryFollowers(post.UserID)

		// loc nguoi theo doi dang online || active
		// onlineFollowers := filterOnlineUsers(followers)

		// cap nhat cache cho timeline service
		// updateTimelineCache(post, onlineFollowers)
	}
}

func updateTimelineCache(post model.Post, onlineFollowers []model.UserInfo) {
	for _, follower := range onlineFollowers {
		key := fmt.Sprintf("timeline:%d", follower.UserID)
		timeline := getCache(key)
		timeline = append(timeline, post)
		setCache(key, timeline)
	}
}
func getCache(key string) []model.Post {
	var timeline []model.Post
	data, _ := global.Rdb.Get(context.Background(), key).Result()
	json.Unmarshal([]byte(data), &timeline)
	return timeline
}

func setCache(key string, timeline []model.Post) {
	data, _ := json.Marshal(timeline)
	global.Rdb.Set(context.Background(), key, data, 0)
}
