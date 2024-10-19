package service

import (
	"context"
	model "go-ecommerce-backend-api/m/v2/internal/models"
)

type (
	// ..interface
	IUserLogin interface {
		Login(ctx context.Context) error
		Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error)
		VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error)
		UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error)
	}

	IUserInfor interface {
		GetInfoByUserId(ctx context.Context) error
		GetAllUser(ctx context.Context) error
	}

	IUserAdmin interface {
		RemoveUser(ctx context.Context) error
		FindOneUser(ctx context.Context) error
	}
)

var (
	localUserAdmin IUserAdmin
	localUserInfo  IUserInfor
	localUserLogin IUserLogin
)

func UserAdmin() IUserAdmin {

	if localUserAdmin == nil {
		panic("implement localUserAdmin notfound")
	}

	return localUserAdmin
}

func InitUserAdmin(i IUserAdmin) {
	localUserAdmin = i
}

func UserInfo() IUserInfor {

	if localUserInfo == nil {
		panic("implement localUserInfo notfound")
	}

	return localUserInfo
}

func InitUserInfo(i IUserInfor) {
	localUserInfo = i
}

func UserLogin() IUserLogin {

	if localUserLogin == nil {
		panic("implement localUserLogin notfound")
	}

	return localUserLogin
}

func InitUserLogin(i IUserLogin) {
	localUserLogin = i
}
