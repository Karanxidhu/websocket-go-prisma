package service

import (
	"context"

	"github.com/karanxidhu/go-websocket/data/request"
)

type FileService interface {
	Save(ctx context.Context, file request.FileResponse)
	GetChat(ctx context.Context, roomName string) (interface{}, error)
}
