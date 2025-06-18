package impl

import (
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/package/utils"
	"go-ecommerce-backend-api/m/v2/response"
	"log"

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
func (s *sTimeline) GetAllPosts(ctx *gin.Context, userId int64) (int, []model.Post, error) {
	cacheKey := fmt.Sprintf("timeline:%d", userId)

	// 1. Kiểm tra cache trước
	if cachedPosts := getCache(cacheKey); cachedPosts != nil {
		fmt.Println("cached post found for key:", cacheKey)
		return response.ErrCodeSuccess, cachedPosts, nil
	}

	// 2. Nếu không có cache, truy vấn DB
	log.Println("Cache miss, querying database for userId:", userId)
	rows, err := s.r.GetAllpost(ctx, sql.NullInt64{Int64: userId, Valid: userId > 0})
	if err != nil {
		return response.ErrCodePostFailed, nil, fmt.Errorf("failed to get posts from DB: %w", err)
	}

	// 3. Chuyển dữ liệu DB sang model
	posts := make([]model.Post, 0, len(rows))
	for _, row := range rows {
		post, err := utils.MapGetAllpostRowToPost(row)
		if err != nil {
			return response.ErrCodePostFailed, nil, fmt.Errorf("failed to map DB row to Post model: %w", err)
		}
		posts = append(posts, post)
	}

	// 4. Lưu cache
	setCache(cacheKey, &posts)

	return response.ErrCodeSuccess, posts, nil
}

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

	// Truy vấn DB lấy row
	row, err := s.r.GetPostById(ctx, id)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("failed to get post from DB: %v", err)
	}

	// Dùng hàm map để chuyển đổi sang model.Post
	post, err := utils.MapGetPostByIdRowToPost(row)
	if err != nil {
		return response.ErrCodePostFailed, model.Post{}, fmt.Errorf("failed to map DB row to Post model: %v", err)
	}

	// Lưu vào cache
	setCache(cacheKey, &[]model.Post{post})

	return response.ErrCodeSuccess, post, nil
}
