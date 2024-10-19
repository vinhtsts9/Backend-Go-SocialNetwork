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
