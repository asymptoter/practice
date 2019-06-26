package auth

import (
	"github.com/asymptoter/practice/base/ctx"
)

type SignupInfo struct {
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
	NickName string `json:"nickName" gorm:"nick_name"`
}

type LoginInfo struct {
	Email    string `json:"email" gorm:"email"`
	Password string `json:"password" gorm:"password"`
}

type Service interface {
	Signup(context ctx.CTX, signupInfo *SignupInfo) error
	Login(context ctx.CTX, loginInfo *LoginInfo) (string, error)
}
