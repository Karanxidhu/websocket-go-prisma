package model

import "time"

type File struct {
	Id         string
	Url        string
	UploadedAt time.Time
	RoomName     string
}
