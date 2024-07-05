package auth

import (
	"github.com/asymptoter/practice-backend/base/ctx"
	"github.com/asymptoter/practice-backend/external/email"
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

func (s *impl) Signup(ctx ctx.CTX, userInfo *models.User) (*models.User, error) {
	// Store user infomation in db
	user, err := s.userStore.Create(ctx, userInfo)
	if err != nil {
		ctx.Error(err)
		return nil, err
	}

	// Send a email to inform the registration succeeded
	signupSuccessMessage := "Registration succeeded!"
	if err := email.Send(ctx, userInfo.Email, signupSuccessMessage); err != nil {
		ctx.Error(err)
		return nil, err
	}

	return user, nil
}
