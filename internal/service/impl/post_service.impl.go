package impl

import (
	"context"
	"encoding/json"
	"errors"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"strconv"
	"time"

	"github.com/segmentio/kafka-go"
)

type sPost struct {
	r *database.Queries
}

func NewPostImpl(r *database.Queries) *sPost {
	return &sPost{
		r: r,
	}
}

// CreatePost tạo một bài viết mới
func (s *sPost) CreatePost(ctx context.Context, input *model.CreatePostInput) (codeRs int, data model.Post, err error) {
	content, err := json.Marshal(input.Content)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Tạo đối tượng CreatePostParams
	createPostParams := database.CreatePostParams{
		Title:   input.Title,
		Content: content,
		UserID:  input.UserID,
	}

	// Gọi hàm tạo bài viết từ repository
	err = s.r.CreatePost(ctx, createPostParams)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Trả về bài viết vừa tạo
	post := model.Post{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    uint32(input.UserID),
		CreatedAt: time.Now().GoString(),
		UpdatedAt: time.Now().GoString(),
	}

	global.KafkaProducer.Topic = "create-post"
	postData, err := json.Marshal(post)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	err = global.KafkaProducer.WriteMessages(ctx, kafka.Message{
		Value: postData,
	})
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	return response.ErrCodeSuccess, post, nil
}

// UpdatePost cập nhật thông tin bài viết
func (s *sPost) UpdatePost(ctx context.Context, postId string, input *model.UpdatePostInput) (codeRs int, data model.Post, err error) {
	id, err := parsePostId(postId)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	content, err := json.Marshal(input.Content)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Tạo đối tượng UpdatePostParams
	updatePostParams := database.UpdatePostParams{
		Title:   input.Title,
		Content: content,
		ID:      id,
	}

	// Gọi hàm cập nhật bài viết từ repository
	err = s.r.UpdatePost(ctx, updatePostParams)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Trả về bài viết đã được cập nhật
	post := model.Post{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    uint32(input.UserID),  // Assuming UserID stays the same
		CreatedAt: time.Now().GoString(), // Assuming created_at doesn't change on update
		UpdatedAt: time.Now().GoString(),
	}
	return response.ErrCodeSuccess, post, nil
}

// DeletePost xóa bài viết theo ID
func (s *sPost) DeletePost(ctx context.Context, postId string) (codeRs int, err error) {
	id, err := parsePostId(postId)
	if err != nil {
		return response.ErrCodePostFailed, err
	}

	// Gọi hàm xóa bài viết từ repository
	err = s.r.DeletePost(ctx, id)
	if err != nil {
		return response.ErrCodePostFailed, err
	}

	return response.ErrCodeSuccess, nil
}

// Helper function để chuyển đổi postId từ string sang uint64
func parsePostId(postId string) (uint64, error) {
	id, err := strconv.ParseUint(postId, 10, 64)
	if err != nil {
		return 0, errors.New("invalid postId")
	}
	return id, nil
}
