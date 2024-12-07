package impl

import (
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/response"

	"github.com/gin-gonic/gin"
)

type sChat struct {
	r *database.Queries
}

func NewsChat(r *database.Queries) *sChat {
	return &sChat{
		r: r,
	}
}

func (s *sChat) CreateRoom(ctx *gin.Context, model *model.CreateRoom) (codeRs int, err error) {
	params := database.CreateRoomChatParams{
		Name:      model.NameRoom,
		IsGroup:   model.IsGroup,
		AdminID:   model.AdminId,
		AvatarUrl: model.AvatarUrl,
	}
	err = s.r.CreateRoomChat(ctx, params)
	if err != nil {
		return response.ErrCodeCreateRoom, err
	}
	return response.ErrCodeSuccess, nil
}
