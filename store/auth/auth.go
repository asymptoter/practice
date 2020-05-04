package auth

import (
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/base/email"
	"github.com/asymptoter/practice-backend/models"
	"github.com/asymptoter/practice-backend/store/user"
)

type Store interface {
	Signup(context ctx.CTX, userInfo *models.User) (*models.User, error)
}

type impl struct {
	userStore user.Store
}

func New(userStore user.Store) Store {
	return &impl{
		userStore: userStore,
	}
}

func (s *impl) Signup(context ctx.CTX, userInfo *models.User) (*models.User, error) {
	// Store user infomation in db
	user, err := s.userStore.Create(context, userInfo)
	if err != nil {
		context.WithField("err", err).Error("Signup failed at userStore.Create")
		return nil, err
	}

	// Send a email to inform the registration succeeded
	signupSuccessMessage := "Registration succeeded!"
	if err := email.Send(context, userInfo.Email, signupSuccessMessage); err != nil {
		context.WithField("err", err).Error("Signup failed at email.Send")
		return nil, err
	}

	return user, nil
}
