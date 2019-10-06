package auth

import (
	"context"
	"log"

	api "github.com/asymptoter/geochellange/apis/auth"
)

type AuthHandler struct{}

func (h *AuthHandler) Signup(ctx context.Context, req *api.SignupRequest) (resp *api.SignupReply, err error) {
	log.Printf("receive client request, client send:%s\n", req)
	return &api.SignupReply{}, nil
}
