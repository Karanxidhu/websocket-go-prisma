package model

import "time"

type MediaFile struct {
	Id         string
	Url        string
	Type       string
	UploadedAt time.Time
	RoomId     string
}
