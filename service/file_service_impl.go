package service

import (
	"context"
	"fmt"

	"github.com/karanxidhu/go-websocket/data/request"
	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/repository"
)

type FileServiceImpl struct {
	FileRepository repository.FileRepository
}

func NewFileServiceImpl(fileRepository repository.FileRepository) *FileServiceImpl {
	return &FileServiceImpl{
		FileRepository: fileRepository,
	}
}

func (p *FileServiceImpl) Save(ctx context.Context, file request.FileResponse) {
	fileData := model.File{
		Url:        file.Url,
		RoomName:   file.RoomName,
	}
	fmt.Println(fileData)
	fmt.Println("reached service")
	p.FileRepository.Save(ctx, fileData)
}

func (p *FileServiceImpl) GetChat(ctx context.Context, roomName string) (interface{}, error) {
	return p.FileRepository.GetChat(ctx, roomName)
}