package auth

import (
	"github.com/asymptoter/practice/base/ctx"
	"github.com/jinzhu/gorm"
	"github.com/sirupsen/logrus"
)

type impl struct {
	mysql *gorm.DB
}

func NewService(db *gorm.DB) Service {
	return &impl{
		mysql: db,
	}
}

func (s *impl) Signup(context ctx.CTX, signupInfo *SignupInfo) error {
	if err := s.mysql.Create(signupInfo).Error; err != nil {
		context.WithFields(logrus.Fields{
			"email":    signupInfo.Email,
			"password": signupInfo.Password,
		}).Error("mysql.Create failed")
		return err
	}
	return nil
}

func (s *impl) Login(context ctx.CTX, loginInfo *LoginInfo) (string, error) {

	return "", nil
}
