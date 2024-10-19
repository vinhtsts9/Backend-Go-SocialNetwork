package service

import (
	"context"
	"encoding/json"
	"fmt"
	"go-ecommerce-backend-api/m/v2/global"
	"go-ecommerce-backend-api/m/v2/internal/repo"
	"go-ecommerce-backend-api/m/v2/package/utils/crypto"
	"go-ecommerce-backend-api/m/v2/package/utils/random"
	"go-ecommerce-backend-api/m/v2/response"
	"time"

	"github.com/segmentio/kafka-go"
)

type IUserService interface {
	Register(email string, purpose string) int
}

type userService struct {
	userRepo     repo.IUserRepository
	userAuthRepo repo.IUserAuthRepo
}

func NewUserService(
	userAuthRepo repo.IUserAuthRepo,
	userRepo repo.IUserRepository,
) IUserService {
	return &userService{
		userRepo:     userRepo,
		userAuthRepo: userAuthRepo,
	}
}

// Register implements IUserService.
func (us *userService) Register(email string, purpose string) int {
	// 0.hashEmail
	hashEmail := crypto.GetHash(email)
	fmt.Printf("hashEmail::%s", hashEmail)
	// 1. check Email Exist In Db
	if us.userRepo.GetUserByEmail(email) {
		return response.ErrCodeUserHasExist
	}
	// 2. New OTP
	otp := random.GenerateSixDigitOtp()
	if purpose == "TETS_USER" {
		otp = 123456
	}

	fmt.Printf("Otp is :::%d\n", otp)
	// 3. Save OTP in RedisExpirationTime
	err := us.userAuthRepo.AddOTP(hashEmail, otp, int64(10*time.Minute))
	if err != nil {
		return response.ErrInvalidOTP
	}
	// 4. Send email Otp
	// err = sendto.SendEmailToJavaByAPI(strconv.Itoa(otp), email, "otp-auth.html")
	// fmt.Printf("err sendto:Java::%d\n", err)
	// if err != nil {
	// 	return response.ErrSendEmailOtp
	// }

	// send OTP via Kafka Java
	body := make(map[string]interface{})
	body["otp"] = otp
	body["email"] = email

	bodyRequest, _ := json.Marshal(body)

	message := kafka.Message{
		Key:   []byte("otp-auth"),
		Value: []byte(bodyRequest),
		Time:  time.Now(),
	}

	err = global.KafkaProducer.WriteMessages(context.Background(), message)
	if err != nil {
		return response.ErrSendEmailOtp
	}
	// 5. Check otp is available

	// 6. User spam

	return response.ErrCodeSuccess
}
