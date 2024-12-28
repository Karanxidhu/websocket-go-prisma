package repository

import (
	"context"
	"fmt"

	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/prisma/db"
)

type FileRepositoryImpl struct {
	Db *db.PrismaClient
}

func NewFileRepository(db *db.PrismaClient) *FileRepositoryImpl {
	return &FileRepositoryImpl{
		Db: db,
	}
}

func (p *FileRepositoryImpl) Save(ctx context.Context, file model.File) error {
	fmt.Println("saving file")
	// Find the room by RoomName
	res, err := p.Db.Room.FindFirst(db.Room.ID.Equals(file.RoomName)).Exec(ctx)
	// if err != nil {
	// 	// Return meaningful error
		fmt.Println(err)
	// 	return fmt.Errorf("failed to find room: %w", err)
	// }
	
	// If no room is found, create one
	if res == nil {
		res, err = p.Db.Room.CreateOne(
			db.Room.Name.Set(file.RoomName),
			).Exec(ctx)
			if err != nil {
			fmt.Println("e2")
			return fmt.Errorf("failed to create room: %w", err)
		}
	}
	
	// Create a MediaFile linked to the Room
	result, err := p.Db.MediaFile.CreateOne(
		db.MediaFile.URL.Set(file.Url),
		db.MediaFile.Room.Link(db.Room.ID.Equals(res.ID)),
		).Exec(ctx)
		if err != nil {
		fmt.Println("e3")
		return fmt.Errorf("failed to create media file: %w", err)
	}

	// Log the success (or use a logger)
	fmt.Println("Media file created successfully:", result)
	return nil
}

func SaveFile(ctx context.Context, file model.File, p *db.PrismaClient) error {
	fmt.Println("saving file")
	// Find the room by RoomName
	res, err := p.Room.FindFirst(db.Room.ID.Equals(file.RoomName)).Exec(ctx)
	// if err != nil {
	// 	// Return meaningful error
		fmt.Println(err)
	// 	return fmt.Errorf("failed to find room: %w", err)
	// }
	
	// If no room is found, create one
	if res == nil {
		res, err = p.Room.CreateOne(
			db.Room.Name.Set(file.RoomName),
			).Exec(ctx)
			if err != nil {
			fmt.Println("e2")
			return fmt.Errorf("failed to create room: %w", err)
		}
	}
	
	// Create a MediaFile linked to the Room
	result, err := p.MediaFile.CreateOne(
		db.MediaFile.URL.Set(file.Url),
		db.MediaFile.Room.Link(db.Room.ID.Equals(res.ID)),
		).Exec(ctx)
		if err != nil {
		fmt.Println("e3")
		return fmt.Errorf("failed to create media file: %w", err)
	}

	// Log the success (or use a logger)
	fmt.Println("Media file created successfully:", result)
	return nil
}
