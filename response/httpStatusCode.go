package response

const (
	ErrCodeSuccess      = 20001
	ErrCodeParamInvalid = 20003

	ErrInvalidToken = 30001
	ErrInvalidOTP   = 30002
	ErrSendEmailOtp = 30003

	ErrCodeUserHasExist = 50001
	// Err login
	ErrCodeOtpNotExist      = 60009
	ErrCodeUserOtpNotExists = 60008
)

var msg = map[int]string{
	ErrCodeSuccess:          "Success",
	ErrCodeParamInvalid:     "Email is invalid",
	ErrInvalidToken:         "Invalid token",
	ErrCodeUserHasExist:     "User has exist",
	ErrInvalidOTP:           "OTP error",
	ErrSendEmailOtp:         "Failed to send email OTP",
	ErrCodeOtpNotExist:      "OTP exists but not registed",
	ErrCodeUserOtpNotExists: "ErrCodeUserOtpNotExists",
}
