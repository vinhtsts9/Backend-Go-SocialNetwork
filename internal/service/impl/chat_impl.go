package impl

import (
	"context"
	"database/sql"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
	"go-ecommerce-backend-api/m/v2/response"
	"time"

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

func (s *sChat) GetRoomChatByUserId(ctx *gin.Context, userId uint64) (codeRs int, rs []model.CreateRoom, err error) {
	Rows, err := s.r.GetRoomByUserId(ctx, userId)
	if err != nil {
		return response.ErrCodeCreateRoom, nil, err
	}
	var rooms []model.CreateRoom
	for _, Row := range Rows {
		room := model.CreateRoom{
			Id:        int32(Row.ID),
			NameRoom:  Row.Name,
			IsGroup:   Row.IsGroup,
			AdminId:   Row.AdminID,
			AvatarUrl: Row.AvatarUrl,
		}
		rooms = append(rooms, room)
	}
	return response.ErrCodeSuccess, rooms, nil

}
func (s *sChat) GetChatHistory(ctx *gin.Context, roomId int) (codeRs int, err error, rs []model.ModelChat) {

	rows, err := s.r.GetChatHistory(ctx, sql.NullInt32{Int32: int32(roomId), Valid: true})
	if err != nil {
		return response.ErrCodeGetMessage, err, nil
	}
	var chatHistory []model.ModelChat
	for _, row := range rows {
		chat := model.ModelChat{
			UserNickname:   row.UserNickname,
			MessageContext: row.MessageContext,
			MessageType:    model.MessagesMessageType(row.MessageType),
			IsPinned:       row.IsPinned,
			CreatedAt:      row.CreatedAt,
		}
		chatHistory = append(chatHistory, chat)
	}
	return response.ErrCodeSuccess, nil, chatHistory
}
func (s *sChat) SetChatHistory(ctx *gin.Context, model *model.ModelChat) {
	dbCtx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	params := database.SetChatHistoryParams{
		UserNickname:   model.UserNickname,
		MessageContext: model.MessageContext,
		MessageType:    database.MessagesMessageType(model.MessageType),
		RoomID:         model.RoomId,
	}
	err := s.r.SetChatHistory(dbCtx, params)
	if err != nil {
		global.Logger.Sugar().Error(err)
		response.ErrorResponse(ctx, response.ErrCodeGetMessage, err.Error())
		return
	}
}

func (s *sChat) GetUserNickName(ctx *gin.Context) (codeRs int, rs string, err error) {
	userInfo := auth.GetUserInfoFromContext(ctx)
	if userInfo == (model.UserInfo{}) {
		return response.ErrCodeGetMessage, "", err
	}
	return response.ErrCodeSuccess, userInfo.UserNickname.String, nil
}

func (s *sChat) DeleteMemberFromGroup(ctx *gin.Context, userid uint64, roomId int64) (codeRs int, Rs bool, err error) {
	params := database.DeleteMemberFromRoomChatParams{
		UserID: userid,
		RoomID: roomId,
	}
	err = s.r.DeleteMemberFromRoomChat(ctx, params)
	if err != nil {
		return response.ErrCodeGetMessage, false, err
	}
	return response.ErrCodeSuccess, true, nil

}
func (s *sChat) GetMemberGroup(ctx *gin.Context, roomId int64) (codeRs int, Rs []model.UserSearch, err error) {
	rows, err := s.r.GetMemberGroup(ctx, int64(roomId))
	if err != nil {
		return response.ErrCodeGetMessage, nil, err
	}
	var results []model.UserSearch
	for _, row := range rows {
		result := model.UserSearch{
			UserNickname: row.UserNickname.String,
			UserAvatar:   row.UserAvatar.String,
		}
		results = append(results, result)
	}
	return response.ErrCodeSuccess, results, nil
}
