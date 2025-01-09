package impl

import (
	"context"
	"database/sql"
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
	var rightValue int32
	global.Logger.Sugar().Info(Comment)
	parent, err := s.r.GetCommentByID(reqCtx, Comment.CommentParentId)
	if err != nil {
		global.Logger.Sugar().Info(0)
		rightValue = 0
	} else {
		rightValue = parent.CommentRight
		global.Logger.Sugar().Info(1)
	}
	if rightValue > 0 {

		// Sử dụng WaitGroup để chờ đợi các goroutine hoàn thành
		var wg sync.WaitGroup
		wg.Add(3)
		go func() {
			defer wg.Done()
			if err := s.r.AddReplyCommentParent(ctx, Comment.CommentParentId); err != nil {
				global.Logger.Sugar().Error("update reply for parent failed")
			}
		}()
		go func() {
			defer wg.Done() // Đảm bảo goroutine này hoàn thành khi hết việc
			params := database.UpdateCommentRightCreateParams{
				PostID:       Comment.PostId,
				CommentRight: rightValue,
			}
			if err := s.r.UpdateCommentRightCreate(reqCtx, params); err != nil {
				global.Logger.Sugar().Error("update cmRight failed")
			}
		}()

		go func() {
			defer wg.Done()
			params := database.UpdateCommentLeftCreateParams{
				PostID:      Comment.PostId,
				CommentLeft: rightValue,
			}
			if err := s.r.UpdateCommentLeftCreate(reqCtx, params); err != nil {
				global.Logger.Sugar().Error("update cmLeft failed")
			}
		}()

		// Chờ đợi các goroutines hoàn thành
		wg.Wait()

	} else {
		MaxRightValue, err := s.r.GetMaxRightComment(reqCtx, Comment.PostId)
		if err != nil {
			if err == sql.ErrNoRows {
				rightValue = 1
			} else {
				return 1, model.ListCommentOutput{}, err
			}
		} else {
			rightValue = MaxRightValue + 1
		}

	}
	global.Logger.Sugar().Info(2)

	// Sử dụng goroutine cho CreateComment nhưng đợi kết quả hoàn thành (thông qua WaitGroup)

	dbs := database.CreateCommentParams{
		PostID:         Comment.PostId,
		UserID:         userId,
		UserNickname:   Comment.UserNickname,
		CommentContent: Comment.CommentContent,
		CommentLeft:    rightValue,
		CommentRight:   rightValue + 1,
		CommentParent: sql.NullInt32{
			Int32: Comment.CommentParentId,
			Valid: Comment.CommentParentId != 0,
		},
		Isdeleted: sql.NullBool{
			Bool:  false,
			Valid: true,
		},
	}
	result, err := s.r.CreateComment(reqCtx, dbs)
	if err != nil {
		global.Logger.Sugar().Error("Create comment failed", err)
	}

	// Lấy ID của bản ghi vừa được chèn và kiểm tra lỗi
	lastInsertId, err := result.LastInsertId()
	if err != nil {
		global.Logger.Sugar().Error("Failed to get last insert ID", err)
		return
	}

	// Gọi hàm GetCommentByLastInsertId với ID vừa lấy được
	CommentRs, err := s.r.GetCommentByLastInsertId(reqCtx, int32(lastInsertId))
	if err != nil {
		global.Logger.Sugar().Error("Get comment by last insert ID failed", err)
		return
	}
	Rs := model.ListCommentOutput{
		Id:              CommentRs.ID,
		PostId:          CommentRs.PostID,
		UserNickname:    CommentRs.UserNickname,
		CommentContent:  CommentRs.CommentContent,
		CommentParentId: CommentRs.CommentParent.Int32,
		Isdeleted:       CommentRs.Isdeleted.Bool,
		CreatedAt:       CommentRs.CreatedAt.Time.Format(time.RFC3339),
	}

	return response.ErrCodeSuccess, Rs, nil
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
