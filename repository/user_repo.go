package repository

import (
	"context"

	"github.com/karanxidhu/go-websocket/model"
)

type UserRepository interface {
	Save(ctx context.Context, user model.User) (string, error)
	Update(ctx context.Context, user model.User)
	Delete(ctx context.Context, userId string)
	FindById(ctx context.Context, userId string) (model.User, error)
	FindAll(ctx context.Context) []model.User
}
