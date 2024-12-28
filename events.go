package main

import "encoding/json"

type Event struct{
	Type string `json:"type"`
	Payload json.RawMessage `json:"payload"`
}

type EventHandler func(event Event, client *Client) error

const (
	JoinEvent = "join"
	LeaveEvent = "leave"
	MessageEvent = "message"
	File = "file"
	WelcomeEvent = "welcome"
)

type Join struct{
	Name string `json:"name"`
}

type Leave struct{
	Name string `json:"name"`
}

type Message struct{
	Name string `json:"name"`
	Message string `json:"message"`
}

type FileUpload struct{
	Name string `json:"name"`
	File string `json:"file"`
}