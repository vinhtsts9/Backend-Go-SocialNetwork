package impl

import (
	"context"
	"database/sql"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	consts "go-ecommerce-backend-api/m/v2/internal/const"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/package/utils"
	"go-ecommerce-backend-api/m/v2/package/utils/crypto"
	"go-ecommerce-backend-api/m/v2/package/utils/random"
	"go-ecommerce-backend-api/m/v2/package/utils/sendto"
	"go-ecommerce-backend-api/m/v2/response"
	"log"
	"strconv"
	"strings"
	"time"

	"github.com/redis/go-redis/v9"
)

type sUserLogin struct {
	r *database.Queries
}

func NewUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

func (s *sUserLogin) Login(ctx context.Context) error {
	return nil
}

func (s *sUserLogin) Register(ctx context.Context, in *model.RegisterInput) (codeResult int, err error) {
	// 1. hash email
	fmt.Printf("VerifyKey: %s\n", in.VerifyKey)
	fmt.Printf("VerifyType: %d\n", in.VerifyType)
	fmt.Printf("VerifyPurpose: %s\n", in.VerifyPurpose)
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))
	fmt.Printf("hashKey: %s\n", hashKey)

	// 2. check user exists in user base
	userFound, err := s.r.CheckUserBaseExists(ctx, in.VerifyKey)
	if err != nil {
		return response.ErrCodeUserHasExist, err
	}
	if userFound > 0 {
		return response.ErrCodeUserHasExist, fmt.Errorf("user has already registed")
	}
	// 3. create otp
	userKey := utils.GetUserKey(hashKey)
	otpFound, err := global.Rdb.Get(ctx, userKey).Result()
	switch {
	case err == redis.Nil:
		fmt.Println("Key doesnt exist")
	case err != redis.Nil:
		fmt.Println("Get failed:", err)
		return response.ErrInvalidOTP, err
	case otpFound != "":
		return response.ErrCodeOtpNotExist, fmt.Errorf("")
	}
	// 4. Generate OTP
	otpNew := random.GenerateSixDigitOtp()
	if in.VerifyPurpose == "TEST_USER" {
		otpNew = 123456
	}
	fmt.Printf("OTP is::%d\n", otpNew)
	// 5. Save OTP in Redis with expiratonTime
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}
	// 6. Send OTP
	switch in.VerifyType {
	case consts.EMAIL:
		err = sendto.SendTextEmail([]string{in.VerifyKey}, consts.HOST_EMAIL, strconv.Itoa(otpNew))
		if err != nil {
			return response.ErrSendEmailOtp, err
		}
		result, err := s.r.InsertOTPVerify(ctx, database.InsertOTPVerifyParams{
			VerifyOtp:     strconv.Itoa(otpNew),
			VerifyType:    sql.NullInt32{Int32: 1, Valid: true},
			VerifyKey:     in.VerifyKey,
			VerifyKeyHash: hashKey,
		})

		if err != nil {
			return response.ErrSendEmailOtp, err
		}
		// 8. getlastId
		lastIdVerifyIdUser, err := result.LastInsertId()
		if err != nil {
			return response.ErrSendEmailOtp, err
		}
		log.Println("lastIdVerifyUser", lastIdVerifyIdUser)

		return response.ErrCodeSuccess, nil
	case consts.MOBILE:
		return response.ErrCodeSuccess, nil
	}
	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) VerifyOTP(ctx context.Context, in *model.VerifyInput) (out model.VerifyOTPOutput, err error) {
	// logic
	hashKey := crypto.GetHash(strings.ToLower(in.VerifyKey))

	// get otp
	otpFound, err := global.Rdb.Get(ctx, utils.GetUserKey(hashKey)).Result()
	if err != nil {
		return out, err
	}
	if in.VerifyCode != otpFound {
		return out, fmt.Errorf("OTP not match")
	}
	infoOTP, err := s.r.GetInfoOTP(ctx, hashKey)
	if err != nil {
		return out, err
	}
	// update status verified
	err = s.r.UpdateUserVerificationStatus(ctx, hashKey)
	if err != nil {
		return out, err
	}

	out.Token = infoOTP.VerifyKeyHash
	out.Message = "success"

	return out, err
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, token string, password string) (userId int, err error) {
	infoOTP, err := s.r.GetInfoOTP(ctx, token)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	// check isVerified ok
	if infoOTP.IsVerified.Int32 == 0 {
		return response.ErrCodeUserOtpNotExists, fmt.Errorf("user OTP not verified")
	}
	// update user_base table
	userBase := database.AddUserBaseParams{}
	userBase.UserAccount = infoOTP.VerifyKey
	userSalt, err := crypto.GenerateSalt(16)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	userBase.UserSalt = userSalt
	userBase.UserPassword = crypto.HassPassword(password, userSalt)
	// add userBase to userBase table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	// add user_id to user_info table
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:             uint64(user_id),
		UserAccount:        infoOTP.VerifyKey,
		UserNickname:       sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserAvatar:         sql.NullString{String: "", Valid: true},
		UserState:          1,
		UserMobile:         sql.NullString{String: "", Valid: true},
		UserGender:         sql.NullInt16{Int16: 0, Valid: true},
		UserBirthday:       sql.NullTime{Time: time.Time{}, Valid: false},
		UserEmail:          sql.NullString{String: infoOTP.VerifyKey, Valid: true},
		UserIsAuthencation: 1,
	})
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	user_id, err = newUserInfo.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}

	return int(user_id), nil
}
