package impl

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	consts "go-ecommerce-backend-api/m/v2/internal/const"
	"go-ecommerce-backend-api/m/v2/internal/database"
	model "go-ecommerce-backend-api/m/v2/internal/models"
	"go-ecommerce-backend-api/m/v2/package/utils"
	"go-ecommerce-backend-api/m/v2/package/utils/auth"
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

func NemUserLoginImpl(r *database.Queries) *sUserLogin {
	return &sUserLogin{
		r: r,
	}
}

// -- Two Factor authen
func (s *sUserLogin) IsTwoFactorEnabled(ctx context.Context, userId int) (codeResult int, rs bool, err error) {
	return 200, true, nil
}

func (s *sUserLogin) SetupTwoFactorAuth(ctx context.Context, in *model.SetupTwoFactorAuthInput) (codeResult int, err error) {
	// 1.Check istwoFactorEnable
	isTwoFactorAuth, err := s.r.IsTwoFactorEnabled(ctx, in.UserId)
	if err != nil {
		return response.ErrCodeTwoFactorSetupFailed, err
	}
	if isTwoFactorAuth > 0 {
		return response.ErrCodeTwoFactorSetupFailed, err
	}
	// 2. Create new type Authe
	err = s.r.EnableTwoFactorTypeEmail(ctx, database.EnableTwoFactorTypeEmailParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.UserTwoFactorTwoFactorAuthTypeEMAIL,
		TwoFactorEmail:    sql.NullString{String: in.TwoFactorEmail, Valid: true},
	})
	if err != nil {
		return response.ErrCodeTwoFactorSetupFailed, err
	}
	//3 . send otp to in.TwoFactorEmail
	keyUserTwoFactor := crypto.GetHash("2fa" + strconv.Itoa(int(in.UserId)))
	go global.Rdb.Set(ctx, keyUserTwoFactor, "123456", time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()

	return response.ErrCodeSuccess, nil
}

func (s *sUserLogin) VerifyTwoFactorAuth(ctx context.Context, in *model.TwoFactorVerificationInput) (codeResult int, err error) {
	// 1. Lấy mã OTP từ Redis
	keyUserTwoFactor := crypto.GetHash("2fa" + strconv.Itoa(int(in.UserId)))
	storedOtp, err := global.Rdb.Get(ctx, keyUserTwoFactor).Result()
	if err == redis.Nil {
		return response.ErrCodeOtpNotExist, fmt.Errorf("OTP does not exist")
	} else if err != nil {
		return response.ErrCodeOtpNotExist, err
	}

	// 2. So sánh mã OTP
	if in.TwoFactorCode != storedOtp {
		return response.ErrCodeOtpMismatch, fmt.Errorf("OTP does not match")
	}

	// 3. Xóa mã OTP sau khi xác thực thành công
	err = global.Rdb.Del(ctx, keyUserTwoFactor).Err()
	if err != nil {
		return response.ErrCodeOtpDeleteFailed, err
	}

	// 4. Cập nhật trạng thái xác thực hai yếu tố
	params := database.UpdateTwoFactorStatusParams{
		UserID:            in.UserId,
		TwoFactorAuthType: database.UserTwoFactorTwoFactorAuthTypeEMAIL,
	}
	err = s.r.UpdateTwoFactorStatus(ctx, params)
	if err != nil {
		return response.ErrCodeTwoFactorUpdateFailed, err
	}

	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) Logout(ctx context.Context, in *model.LogoutInput) (codeRs int, err error) {
	claims := auth.CheckAuth(in.TokenString)
	if err != nil {
		return 1, fmt.Errorf("invalid token: %v", err)
	}
	err = global.Rdb.Del(ctx, claims.Subject).Err()
	if err != nil {
		return 2, fmt.Errorf("Failed to delete session from Redis: %v", err)

	}
	userAccount := auth.GetUserInfoFromToken(in.TokenString)
	err = s.r.LogoutUserBase(ctx, userAccount.UserAccount)
	if err != nil {
		return response.ErrCodeNotFound, err
	}

	return response.ErrCodeSuccess, nil
}
func (s *sUserLogin) Login(ctx context.Context, in *model.LoginInput) (codeResult int, out model.LoginOutPut, err error) {
	// logic login
	fmt.Println(in.UserAccount)
	userBase, err := s.r.GetOneUserInfo(ctx, in.UserAccount)
	if err != nil {
		return 1, out, err
	}
	// 2. check password
	if !crypto.MatchingPassword(userBase.UserPassword, in.UserPassword, userBase.UserSalt) {
		return 2, out, fmt.Errorf("does not password")
	}
	// 3. Check two factor authentication

	// 4. Update password time
	go s.r.LoginUserBase(ctx, database.LoginUserBaseParams{
		UserLoginIp: sql.NullString{String: "127.0.0.1", Valid: true},
		UserAccount: in.UserAccount,
	})
	// 5. Create uuid User
	subToken := utils.GenerateCliTokenUUID(int(userBase.UserID))

	// 6. get user_info table
	infoUser, err := s.r.GetUser(ctx, uint64(userBase.UserID))
	if err != nil {
		return 3, out, nil
	}
	// convert to json
	infoUserJson, err := json.Marshal(infoUser)
	if err != nil {
		return 4, out, fmt.Errorf("convert to json failed: %v", err)
	}
	// 7. give infoUserJson to redis with key = subToken
	err = global.Rdb.Set(ctx, subToken, infoUserJson, time.Duration(consts.TIME_OTP_REGISTER)*time.Hour).Err()
	if err != nil {
		return 5, out, err
	}
	// 8. create token
	out.Token, err = auth.CreateToken(subToken)
	if err != nil {
		return 6, out, err
	}
	return 200, out, nil
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
	case err != nil:
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
	err = global.Rdb.SetEx(ctx, userKey, strconv.Itoa(otpNew), time.Duration(consts.TIME_OTP_REGISTER)*time.Minute).Err()
	if err != nil {
		return response.ErrInvalidOTP, err
	}

	fmt.Printf("OTP is::%d\n", otpNew)
	// 5. Save OTP in Redis with expiratonTime

	// 6. Send OTP
	switch in.VerifyType {
	case consts.EMAIL:
		err := sendto.SendTextEmail([]string{in.VerifyKey}, "Vinhtiensinh17@gmail.com", strconv.Itoa(otpNew))
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
	global.Logger.Sugar().Info(out.Token)
	return out, err
}

func (s *sUserLogin) UpdatePasswordRegister(ctx context.Context, model *model.UpdatePasswordRegisterInput) (userId int, err error) {
	infoOTP, err := s.r.GetInfoOTP(ctx, model.UserToken)
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
	userBase.UserPassword = crypto.HassPassword(model.UserPassword, userSalt)
	// add userBase to userBase table
	newUserBase, err := s.r.AddUserBase(ctx, userBase)
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	user_id, err := newUserBase.LastInsertId()
	if err != nil {
		return response.ErrCodeUserOtpNotExists, err
	}
	global.Logger.Sugar().Info(model)
	// add user_id to user_info table
	newUserInfo, err := s.r.AddUserHaveUserId(ctx, database.AddUserHaveUserIdParams{
		UserID:             uint64(user_id),
		UserAccount:        infoOTP.VerifyKey,
		UserNickname:       sql.NullString{String: model.UserNickname, Valid: true},
		UserAvatar:         sql.NullString{String: model.UserAvatar, Valid: true},
		UserState:          1,
		UserMobile:         sql.NullString{String: model.UserMobile, Valid: true},
		UserGender:         sql.NullInt16{Int16: model.UserGender, Valid: true},
		UserBirthday:       sql.NullTime{Time: model.UserBirthday, Valid: false},
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
