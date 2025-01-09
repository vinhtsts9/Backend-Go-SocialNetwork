package impl

import (
	"fmt"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"time"

	"github.com/gin-gonic/gin"
)

type sTimeline struct {
	r *database.Queries
}

// NewTimelineImpl khởi tạo mới đối tượng sTimeline
func NewTimelineImpl(r *database.Queries) *sTimeline {
	return &sTimeline{
		r: r,
	}
}

// GetAllPosts lấy danh sách tất cả bài viết
func (s *sTimeline) GetAllPosts(ctx *gin.Context, userId int64) (codeRs int, data []model.Post, err error) {
	// // Kiểm tra cache trước khi truy vấn database
	cacheKey := fmt.Sprintf("timeline:%d", userId)
	// cachedPosts := getCache(cacheKey)
	// if err == nil && cachedPosts != nil {
	// 	return response.ErrCodeSuccess, cachedPosts, nil
	// }

	// Nếu không có trong cache, truy vấn từ cơ sở dữ liệu
	rows, err := s.r.GetAllpost(ctx)
	if err != nil {
		return response.ErrCodePostFailed, nil, fmt.Errorf("failed to get posts from DB: %v", err)
	}

	var posts []model.Post
	for _, row := range rows {
		posts = append(posts, model.Post{
			ID:           uint32(row.ID),
			UserNickname: row.UserNickname,
			Title:        row.Title,
			ImagePaths:   &row.ImagePaths,
			CreatedAt:    row.CreatedAt.Time.Format(time.RFC3339),
			UpdatedAt:    row.UpdatedAt.Time.Format(time.RFC3339),
		})
	}

	// Lưu bài viết vào cache để sử dụng lần sau
	setCache(cacheKey, posts)
	return response.ErrCodeSuccess, posts, nil
}

// GetPostById lấy bài viết theo ID
func (s *sTimeline) GetPostById(ctx *gin.Context, postId string) (codeRs int, data model.Post, err error) {
	id, err := parsePostId(postId)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("invalid post ID: %v", err)
	}

	// Kiểm tra cache trước khi truy vấn database
	cacheKey := fmt.Sprintf("post:%d", id)
	cachedPost := getCache(cacheKey)
	if err == nil && cachedPost != nil {
		return response.ErrCodeSuccess, cachedPost[0], nil
	}

	// Nếu không có trong cache, truy vấn từ cơ sở dữ liệu
	row, err := s.r.GetPostById(ctx, id)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("failed to get post from DB: %v", err)
	}

	post := model.Post{
		ID:           uint32(row.ID),
		UserNickname: row.UserNickname,
		Title:        row.Title,
		ImagePaths:   &row.ImagePaths,
		CreatedAt:    row.CreatedAt.Time.Format(time.RFC3339),
		UpdatedAt:    row.UpdatedAt.Time.Format(time.RFC3339),
	}

	// Lưu bài viết vào cache để sử dụng lần sau
	setCache(cacheKey, []model.Post{post})

	return response.ErrCodeSuccess, post, nil
}
