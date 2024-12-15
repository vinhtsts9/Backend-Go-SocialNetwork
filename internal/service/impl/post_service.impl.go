package impl

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
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
	// Kiểm tra nếu content là một map JSON hợp lệ
	contentMap, ok := input.Content.(map[string]interface{})
	if !ok {
		return response.ErrCodePostFailed, model.Post{}, errors.New("content must be a valid JSON object")
	}

	// Kiểm tra và xử lý image_url hoặc image_path trong content
	var imageUrl string
	if urlStr, exists := contentMap["image_url"].(string); exists {
		// Nếu có image_url, upload ảnh từ URL
		uploadedUrl, err := global.Cloudinary.UploadImageFromURLToCloudinary(urlStr)
		if err != nil {
			return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("failed to process image_url: %v", err)
		}
		imageUrl = uploadedUrl
	} else if imagePath, exists := contentMap["image_path"].(string); exists {
		// Nếu không có image_url nhưng có image_path (ảnh cục bộ)
		uploadedUrl, err := global.Cloudinary.UploadImageToCloudinary(imagePath)
		if err != nil {
			return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("failed to upload image: %v", err)
		}
		imageUrl = uploadedUrl
		delete(contentMap, "image_path") // Xóa image_path nếu không cần thiết
	} else {
		// Nếu không có image_url hoặc image_path, yêu cầu phải có ít nhất một
		return response.ErrCodePostFailed, model.Post{}, errors.New("either image_url or image_path must be provided")
	}

	// Gắn hoặc thay thế image_url trong content nếu có
	if imageUrl != "" {
		contentMap["image_url"] = imageUrl
	}

	// Serialize content để lưu vào database
	contentJSON, err := json.Marshal(input.Content)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, errors.New("failed to serialize post content")
	}

	// Tạo đối tượng CreatePostParams để lưu
	createPostParams := database.CreatePostParams{
		Title:   input.Title,
		Content: contentJSON,
		UserID:  input.UserID,
	}

	// Lưu vào database
	err = s.r.CreatePost(ctx, createPostParams)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, errors.New("failed to create post in database")
	}

	// Trả về bài viết
	post := model.Post{
		Title:     input.Title,
		Content:   input.Content,
		UserID:    uint32(input.UserID),
		CreatedAt: time.Now().Format(time.RFC3339),
		UpdatedAt: time.Now().Format(time.RFC3339),
	}

	// Gửi message vào Kafka
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
