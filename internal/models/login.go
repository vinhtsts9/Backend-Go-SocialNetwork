package model

type RegisterInput struct {
	VerifyKey     string `json:"verify_key"`
	VerifyType    int    `json:"verify_type"`
	VerifyPurpose string `json:"verify_purpose"`
}

type VerifyInput struct {
	VerifyKey  string `json:"verify_key"`
	VerifyCode string `json:"verify_code"`
}

type VerifyOTPOutput struct {
	Token   string `json:"token"`
	UserId  string `json:"userId"`
	Message string `json:"message"`
}

type UpdatePasswordRegisterInput struct {
	UserToken    string `json:"user_token"`
	UserPassword string `json:"user_password"`
}
type LoginInput struct {
	UserAccount  string `json:"user_account"`
	UserPassword string `json:"user_password"`
}
type LoginOutPut struct {
	Token   string `json:"token"`
	Message string `json:"message"`
}

// two factor authentication
type SetupTwoFactorAuthInput struct {
	UserId            uint32 `json:"user_id"`
	TwoFactorAuthType string `json:"two_factor_auth_type"`
	TwoFactorEmail    string `json:"two_factor_email"`
}

type TwoFactorVerificationInput struct {
	UserId        uint32 `json:"user_id"`
	TwoFactorCode string `json:"two_factor_code"`
}
