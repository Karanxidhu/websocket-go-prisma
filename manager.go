package main

import (
	"context"
	"crypto/rand"
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
	"github.com/karanxidhu/go-websocket/config"
	"github.com/karanxidhu/go-websocket/model"
	"github.com/karanxidhu/go-websocket/repository"
)

var (
	rooms    = make(map[string]*Room) // Map to store rooms by name
	upgrader = websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true // Allow any origin for simplicity
		},
	}
)

func RandomString(length int) (string, error) {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	b := make([]byte, length)
	for i := range b {
		num, err := rand.Int(rand.Reader, big.NewInt(int64(len(charset))))
		if err != nil {
			return "", err
		}
		b[i] = charset[num.Int64()]
	}
	return string(b), nil
}

var roomMutex sync.Mutex

func getOrCreateRoom(name string, userName string) *Room {
	roomMutex.Lock()
	defer roomMutex.Unlock()

	if name == "" {
		random, err := RandomString(10)
		if err != nil {
			panic(err)
		}
		name = random
	}

	if room, exists := rooms[name]; exists {
		return room
	}
	room := &Room{
		clients: make(map[*Client]bool),
		name:    name,
	}
	rooms[name] = room
	return room
}

type Room struct {
	clients map[*Client]bool // Active clients in the room
	name    string
	mu      sync.Mutex // Name of the room
}

func (room *Room) broadcast(message Event, sender *Client) {
	room.mu.Lock()
	defer room.mu.Unlock()

	if message.Type == WelcomeEvent {
		sender.egress <- message
	}

	for client := range room.clients {
		if client == sender {
			continue // Skip the sender
		}
		select {
		case client.egress <- message:
		default:
		}
	}
}
func (room *Room) GetParticipants() []string {
	room.mu.Lock()
	defer room.mu.Unlock()

	// Collect the names of all clients in the room
	clientNames := []string{}
	for client := range room.clients {
		clientNames = append(clientNames, client.name)
	}

	return clientNames
}

// func (room *Room) welcome(client *Client) {
// 	fmt.Printf("Welcome to room: %s\n", room.name)
// 	room.mu.Lock()
// 	defer room.mu.Unlock()
// 	client.egress <- Event{
// 		Type:    WelcomeEvent,
// 		Payload: []byte(fmt.Sprintf("Welcome to the room: %s", room.name)),
// 	}
// }

type Manager struct {
	clients ClientList
	sync.RWMutex
	handlers map[string]EventHandler
}

func NewManager() *Manager {
	m := &Manager{
		clients:  make(ClientList),
		handlers: make(map[string]EventHandler),
	}
	m.setupEventHandlers()
	return m
}

func (m *Manager) routeEvent(event Event, client *Client) error {
	if handler, ok := m.handlers[event.Type]; ok {
		if err := handler(event, client); err != nil {
			return err
		}
		return nil
	} else {
		return errors.New("there is no event type ")
	}
}

func (m *Manager) setupEventHandlers() {
	m.handlers[JoinEvent] = HandleJoin
	// m.handlers[LeaveEvent] = HandleLeave
	m.handlers[MessageEvent] = HandleMessage
	m.handlers[File] = HandleFile
}

func HandleFile(event Event, client *Client) error {
	var file struct {
		AuthToken string `json:"authToken"`
		FileLink  string `json:"fileLink"`
	}
	if err := json.Unmarshal(event.Payload, &file); err != nil {
		log.Println("Invalid file event:", err)
		return err
	}
	if file.AuthToken == "" {
		return errors.New("auth token is required")
	}

	type Filer struct {
		FileLink string `json:"fileLink"`
	}

	newFile := Filer{
		FileLink: file.FileLink,
	}

	jsonData, err := json.Marshal(newFile)
	if err != nil {
		fmt.Println("Error:", err)
	}

	newEvent := Event{
		Type:    File,
		Payload: jsonData,
	}

	client.room.broadcast(newEvent, client) // Pass client as sender to avoid echoing back
	return nil
}

func HandleMessage(event Event, client *Client) error {
	var msg struct {
		Message  string `json:"message"`
		RoomName string `json:"roomName"`
		UserName string `json:"userName"`
	}
	if err := json.Unmarshal(event.Payload, &msg); err != nil {
		log.Println("Invalid message event:", err)
		return err
	}

	message := model.File{
		Message:  msg.Message,
		RoomName: msg.RoomName,
		UserName: msg.UserName,
	}

	client.room.broadcast(event, client) // Pass client as sender to avoid echoing back
	err := repository.SaveMsg(context.Background(), message, config.Db)

	if err != nil {
		fmt.Println("Error:", err)
	}

	return nil
}

func HandleJoin(event Event, client *Client) error {
	var join struct {
		Name     string `json:"name"`
		UserName string `json:"username"`
	}
	if err := json.Unmarshal(event.Payload, &join); err != nil {
		log.Println("Invalid join event:", err)
		return err
	}

	room := getOrCreateRoom(join.Name, join.UserName)
	client.room = room

	room.mu.Lock()
	room.clients[client] = true
	room.mu.Unlock()

	var pay struct {
		Name       string   `json:"name"`
		UserName   string   `json:"username"`
		Paticipant []string `json:"paticipant"`
	}
	pay.Name = join.Name
	pay.UserName = join.UserName
	pay.Paticipant = room.GetParticipants()

	joined := Event{}
	joined.Type = JoinEvent
	joined.Payload, _ = json.Marshal(pay)

	log.Println("Client joined room: ", join.Name)
	log.Println(joined)
	client.room.broadcast(joined, client)
	return nil
}

func (m *Manager) servesWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new connection recieved")
	roomName := r.URL.Query().Get("room")
	userName := r.URL.Query().Get("name")

	fmt.Println("user name: ", userName)
	if roomName == "" {
		fmt.Println("Room name is required")
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(ws, m, userName)
	m.addClient(client)

	room := getOrCreateRoom(roomName, userName)
	client.room = room

	room.mu.Lock()
	room.clients[client] = true
	room.mu.Unlock()
	go client.writeMessage()

	payload, _ := json.Marshal(map[string]interface{}{
		"username":     userName,
		"name":         roomName,
		"participants": room.GetParticipants(), // Add the array of participants
	})
	
	// Send the welcome message to the client
	welcomeMsg := Event{
		Type:    JoinEvent,
		Payload: json.RawMessage(payload),
	}
	
	client.send(welcomeMsg)

	// room.welcome(client)
	go client.readMessage()
	room.broadcast(welcomeMsg, client)
	// room.broadcast(Event{Type: WelcomeEvent, Payload: []byte(fmt.Sprintf("Welcome to the room: %s", room.name))}, client)

}

func (m *Manager) addClient(client *Client) {
	fmt.Println("new client added")
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	log.Printf("Removing client: %s", client.name)
	m.Lock()
	defer m.Unlock()
	if _, exists := m.clients[client]; exists {
		delete(m.clients, client)
		client.conn.Close()
	}
}

func (room *Room) removeClient(client *Client) {
	room.mu.Lock()
	defer room.mu.Unlock()
	if _, exists := room.clients[client]; exists {
		delete(room.clients, client)
		log.Println("Client removed from room:", room.name, "user name:", client.name)
	}
	repository.RemoveMember(client.name, room.name, context.Background(), config.Db)
}
