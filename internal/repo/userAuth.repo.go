package repo

import (
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"time"
)

type IUserAuthRepo interface {
	AddOTP(email string, otp int, expirationTime int64) error
}

type UserAuthRepo struct{}

// AddOTP implements IUserAuthRepo.
func (u *UserAuthRepo) AddOTP(email string, otp int, expirationTime int64) error {
	key := fmt.Sprintf("usr:%s:otp", email)
	return global.Rdb.SetEx(ctx, key, otp, time.Duration(expirationTime)).Err()
}

func NewUserAuthRepo() IUserAuthRepo {
	return &UserAuthRepo{}

}
