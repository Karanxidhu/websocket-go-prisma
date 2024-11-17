package model

import "time"

type Room struct {
	Id          string
	Name        string
	CreatedAt   time.Time
	MediaFiles  string
}