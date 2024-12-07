package impl

import (
	"context"
	"database/sql"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"
	"sync"
	"time"
)

type sComment struct {
	r *database.Queries
}

func NewCommenService(r *database.Queries) *sComment {
	return &sComment{
		r: r,
	}
}

func (s *sComment) CreateComment(Comment *model.CommentInput) (codeRs int, err error) {
	parent, err := s.r.GetCommentByID(context.Background(), Comment.CommentParentId)
	if err != nil {
		return response.ErrCodeComment, err
	}

	var rightValue int32
	if parent.CommentRight.Valid {
		rightValue = parent.CommentRight.Int32

		go s.r.UpdateCommentRightCreate(context.Background(), Comment.PostId)
		go s.r.UpdateCommentLeftCreate(context.Background(), Comment.PostId)
	} else {
		MaxRightValue, err := s.r.GetMaxRightComment(context.Background(), Comment.PostId)
		if err != nil {
			return response.ErrCodeComment, err
		}
		if MaxRightValue.Valid {
			rightValue = MaxRightValue.Int32 + 1
		} else {
			rightValue = 1
		}
	}
	dbs := database.CreateCommentParams{
		PostID: Comment.PostId,
		UserID: Comment.UserId,
		CommentContent: sql.NullString{
			String: Comment.CommentContent,
			Valid:  Comment.CommentContent != "",
		},
		CommentLeft:   sql.NullInt32{Int32: rightValue, Valid: true},
		CommentRight:  sql.NullInt32{Int32: rightValue + 1, Valid: true},
		CommentParent: sql.NullInt32{Int32: Comment.CommentParentId, Valid: true},
		Isdeleted:     sql.NullBool{Bool: false},
	}
	go s.r.CreateComment(context.Background(), dbs)
	return response.ErrCodeSuccess, nil
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
			Id:              comment.ID,
			PostId:          comment.PostID,
			CommentContent:  comment.CommentContent.String,
			UserId:          comment.UserID,
			CommentParentId: comment.CommentParent.Int32,
			Isdeleted:       comment.Isdeleted.Bool,
		})
	}
	return response.ErrCodeSuccess, result, nil
}

func (s *sComment) DeleteComment(modelInput *model.DeleteCommentInput) (bool, error) {
	// 1. lay thong tin comment can xoa
	comment, err := s.r.GetCommentByID(context.Background(), modelInput.Id)

	if err != nil {
		return false, err
	}

	left := comment.CommentLeft.Int32
	right := comment.CommentRight.Int32

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
			CommentLeft:  sql.NullInt32{Int32: left, Valid: true},
			CommentRight: sql.NullInt32{Int32: right, Valid: true},
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
			CommentLeft:   sql.NullInt32{Int32: width},
			CommentLeft_2: sql.NullInt32{Int32: right},
		}); err != nil {
			errChan <- err
		}
	}()

	// 5. cap nhat comment right
	wg.Add(1)
	go func() {
		if err := s.r.UpdateCommentRight(ctx, database.UpdateCommentRightParams{
			PostID:         modelInput.PostId,
			CommentRight:   sql.NullInt32{Int32: width},
			CommentRight_2: sql.NullInt32{Int32: right},
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
			return false, err
		}
	}
	return true, nil
}
