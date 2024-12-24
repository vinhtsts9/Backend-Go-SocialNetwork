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
	// Log đầu vào
	global.Logger.Sugar().Info("Model create Post, %s", input)

	// Serialize image_paths để lưu vào database
	contentJSON, err := json.Marshal(input.ImagePaths)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, errors.New("failed to serialize image paths")
	}

	// Tạo đối tượng CreatePostParams để lưu vào DB
	createPostParams := database.CreatePostParams{
		Title:        input.Title,
		ImagePaths:   contentJSON,
		UserNickname: input.UserNickname,
		UserID:       input.UserId,
	}

	// Lưu vào database
	err = s.r.CreatePost(ctx, createPostParams)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, errors.New("failed to create post in database")
	}

	// Tạo đối tượng post để trả về
	post := model.Post{
		Title:        input.Title,
		ImagePaths:   input.ImagePaths, // Đây là danh sách file paths
		UserNickname: input.UserNickname,
		CreatedAt:    time.Now().Format(time.RFC3339),
		UpdatedAt:    time.Now().Format(time.RFC3339),
	}

	// Gửi message vào Kafka (nếu cần)
	err = global.KafkaProducer.Send("create-post", post, 3)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, errors.New("failed to send post data to Kafka")
	}

	return response.ErrCodeSuccess, post, nil
}

// UpdatePost cập nhật thông tin bài viết
func (s *sPost) UpdatePost(ctx context.Context, postId string, input *model.UpdatePostInput) (codeRs int, data model.Post, err error) {
	id, err := parsePostId(postId)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	image_paths, err := json.Marshal(input.ImagePaths)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Tạo đối tượng UpdatePostParams
	updatePostParams := database.UpdatePostParams{
		Title:      input.Title,
		ImagePaths: image_paths,
		ID:         id,
	}

	// Gọi hàm cập nhật bài viết từ repository
	err = s.r.UpdatePost(ctx, updatePostParams)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, err
	}

	// Trả về bài viết đã được cập nhật
	post := model.Post{
		Title:        input.Title,
		ImagePaths:   input.ImagePaths,
		UserNickname: input.UserNickname,    // Assuming UserID stays the same
		CreatedAt:    time.Now().GoString(), // Assuming created_at doesn't change on update
		UpdatedAt:    time.Now().GoString(),
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
