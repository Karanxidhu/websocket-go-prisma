package service

import (
	"context"

	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/data/response"
)

type UserService interface {
	Save(ctx context.Context, user request.UserCreateReq) string
	Update(ctx context.Context, user request.UserUpdateReq)
	Delete(ctx context.Context, userId string)
	FindById(ctx context.Context, userId string) response.UserResponse
	FindAll(ctx context.Context) []response.UserResponse
}
