// internal/repository/tx_comment.go

package database

import (
	"context"
	"database/sql"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type CreateCommentTxParams struct {
	Input        *model.CreateCommentInput
	UserID       uint64
	AfterCreated func(comment model.ListCommentOutput) error
}

type CreateCommentTxResult struct {
	Comment model.ListCommentOutput
}

func (s *SQLStore) CreateCommentTx(ctx context.Context, arg CreateCommentTxParams) (CreateCommentTxResult, error) {
	var result CreateCommentTxResult

	err := s.execTx(ctx, func(q *Queries) error {
		var rightValue int32

		// Lấy comment cha nếu có
		parent, err := q.GetCommentByID(ctx, arg.Input.CommentParentId)
		if err != nil && err != sql.ErrNoRows {
			return err
		}

		if parent.ID != 0 {
			rightValue = parent.CommentRight

			if err := q.AddReplyCommentParent(ctx, arg.Input.CommentParentId); err != nil {
				return err
			}

			if err := q.UpdateCommentRightCreate(ctx, UpdateCommentRightCreateParams{
				PostID:       arg.Input.PostId,
				CommentRight: rightValue,
			}); err != nil {
				return err
			}

			if err := q.UpdateCommentLeftCreate(ctx, UpdateCommentLeftCreateParams{
				PostID:      arg.Input.PostId,
				CommentLeft: rightValue,
			}); err != nil {
				return err
			}
		} else {
			maxRight, err := q.GetMaxRightComment(ctx, arg.Input.PostId)
			if err != nil && err != sql.ErrNoRows {
				return err
			}
			if err == sql.ErrNoRows {
				rightValue = 1
			} else {
				rightValue = maxRight + 1
			}
		}

		commentParams := CreateCommentParams{
			PostID:         arg.Input.PostId,
			UserID:         arg.UserID,
			UserNickname:   arg.Input.UserNickname,
			CommentContent: arg.Input.CommentContent,
			CommentLeft:    rightValue,
			CommentRight:   rightValue + 1,
			CommentParent: sql.NullInt32{
				Int32: arg.Input.CommentParentId,
				Valid: arg.Input.CommentParentId != 0,
			},
			Isdeleted: sql.NullBool{
				Bool:  false,
				Valid: true,
			},
		}

		resultRaw, err := q.CreateComment(ctx, commentParams)
		if err != nil {
			return err
		}

		commentId, err := resultRaw.LastInsertId()
		if err != nil {
			return err
		}
		comment, err := q.GetCommentByLastInsertId(ctx, int32(commentId))
		if err != nil {
			return err
		}
		result.Comment = model.ListCommentOutput{
			Id:              comment.ID,
			PostId:          comment.PostID,
			UserNickname:    comment.UserNickname,
			CommentContent:  comment.CommentContent,
			CommentParentId: comment.CommentParent.Int32,
			Isdeleted:       comment.Isdeleted.Bool,
			CreatedAt:       comment.CreatedAt.Time.Format("2006-01-02T15:04:05Z07:00"),
		}

		if arg.AfterCreated != nil {
			return arg.AfterCreated(result.Comment)
		}

		return nil
	})

	return result, err
}
