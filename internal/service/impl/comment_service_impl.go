package impl

import (
	"context"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"sync"
	"time"

	"github.com/gin-gonic/gin"
)

type sComment struct {
	r *database.Queries
}

func NewCommenService(r *database.Queries) *sComment {
	return &sComment{
		r: r,
	}
}
func (s *sComment) CreateComment(ctx *gin.Context, Comment *model.CreateCommentInput, userId uint64) (codeRs int, RS model.ListCommentOutput, err error) {
	reqCtx, cancel := context.WithTimeout(ctx.Request.Context(), 5*time.Second)
	defer cancel()

	params := database.CreateCommentTxParams{
		Input:  Comment,
		UserID: userId,
		AfterCreated: func(comment model.ListCommentOutput) error {
			// Ví dụ log sau khi tạo comment thành công
			global.Logger.Sugar().Infof("Comment created with ID: %d, Content: %s", comment.Id, comment.CommentContent)
			// Có thể gọi các service khác, cập nhật cache, gửi notification...
			return nil
		},
	}

	result, err := global.Store.CreateCommentTx(reqCtx, params)
	if err != nil {
		global.Logger.Sugar().Error("CreateCommentTx failed", err)
		return response.ErrCodeComment, model.ListCommentOutput{}, err
	}

	return response.ErrCodeSuccess, result.Comment, nil
}

func (s *sComment) ListComments(modelInput *model.ListCommentInput) (codeRs int, out []model.ListCommentOutput, err error) {
	query := database.GetCommentByParentIDParams{
		PostID: modelInput.PostId,
		ID:     modelInput.CommentParentId,
		ID_2:   modelInput.CommentParentId,
	}
	comments, err := s.r.GetCommentByParentID(context.Background(), query)
	if err != nil {
		return response.ErrCodeComment, nil, err
	}
	var result []model.ListCommentOutput
	for _, comment := range comments {
		result = append(result, model.ListCommentOutput{
			Id:             comment.ID,
			UserNickname:   comment.UserNickname,
			ReplyCount:     comment.ReplyCount,
			PostId:         comment.PostID,
			CommentContent: comment.CommentContent,
			CreatedAt:      comment.CreatedAt.Time.Format(time.RFC3339),
			Isdeleted:      comment.Isdeleted.Bool,
		})
	}
	return response.ErrCodeSuccess, result, nil
}

func (s *sComment) ListCommentRoot(ctx *gin.Context, postId uint64) (codeRs int, data []model.ListCommentOutput, err error) {
	comments, err := s.r.GetRootComment(ctx, postId)
	if err != nil {
		return response.ErrCodeComment, nil, err
	}
	var result []model.ListCommentOutput
	for _, comment := range comments {
		result = append(result, model.ListCommentOutput{
			Id:             comment.ID,
			UserNickname:   comment.UserNickname,
			PostId:         comment.PostID,
			CommentContent: comment.CommentContent,
			Isdeleted:      comment.Isdeleted.Bool,
			ReplyCount:     comment.ReplyCount,
			CreatedAt:      comment.CreatedAt.Time.Format(time.RFC3339),
		})
	}
	return response.ErrCodeSuccess, result, nil
}

func (s *sComment) DeleteComment(modelInput *model.DeleteCommentInput) (codeRs int, err error, Rs bool) {
	// 1. lay thong tin comment can xoa
	comment, err := s.r.GetCommentByID(context.Background(), modelInput.Id)

	if err != nil {
		return response.ErrCodeComment, err, false
	}

	left := comment.CommentLeft
	right := comment.CommentRight

	width := right - left + 1

	// 2. tao context va waitGroup
	ctx, cancel := context.WithTimeout(context.Background(), time.Second*5)
	defer cancel()

	var wg sync.WaitGroup
	errChan := make(chan error, 3)

	// 3. xoa cac comment trong khoang [left, right]
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.r.DeleteCommentsInRange(ctx, database.DeleteCommentsInRangeParams{
			PostID:       modelInput.PostId,
			CommentLeft:  left,
			CommentRight: right,
		}); err != nil {
			errChan <- err
		}
	}()
	// 4. Cap nhat comment left
	wg.Add(1)
	go func() {
		defer wg.Done()
		if err := s.r.UpdateCommentLeft(ctx, database.UpdateCommentLeftParams{
			PostID:        modelInput.PostId,
			CommentLeft:   width,
			CommentLeft_2: right,
		}); err != nil {
			errChan <- err
		}
	}()

	// 5. cap nhat comment right
	wg.Add(1)
	go func() {
		if err := s.r.UpdateCommentRight(ctx, database.UpdateCommentRightParams{
			PostID:         modelInput.PostId,
			CommentRight:   width,
			CommentRight_2: right,
		}); err != nil {
			errChan <- err
		}
	}()
	// 6. Cho tat ca goroutine hoan thanh
	wg.Wait()
	close(errChan)
	// 7. Gom loi
	for err := range errChan {
		if err != nil {
			return response.ErrCodeComment, err, false
		}
	}
	return response.ErrCodeSuccess, nil, true
}
