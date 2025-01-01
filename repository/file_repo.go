package repository

import (
	"context"

	"github.com/karanxidhu/go-websocket/model"
)

type FileRepository interface {
	Save(ctx context.Context, file model.File) error
	GetChat(ctx context.Context, roomName string) (interface{}, error)
}
