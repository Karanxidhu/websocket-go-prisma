package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"crypto/rand"
	"math/big"
	"net/http"
	"sync"

	"github.com/gorilla/websocket"
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
func getOrCreateRoom(name string) *Room {

	if name == "" {
		random, err := RandomString(10)
			if err != nil {
			panic(err)
			}
		name = random
	}
	fmt.Println("get or create room")
	if room, exists := rooms[name]; exists {
		fmt.Println("room joined with name: ", name)
		return room
	}
	room := &Room{
		clients: make(map[*Client]bool),
		name:    name,
	}
	rooms[name] = room
	fmt.Println("room created with name: ", name)
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
}

func HandleMessage(event Event, client *Client) error {
	var msg struct {
		Message string `json:"message"`
	}
	if err := json.Unmarshal(event.Payload, &msg); err != nil {
		log.Println("Invalid message event:", err)
		return err
	}

	client.room.broadcast(event, client) // Pass client as sender to avoid echoing back
	return nil
}


func HandleJoin(event Event, client *Client) error {
	var join struct {
		Name string `json:"name"`
	}
	if err := json.Unmarshal(event.Payload, &join); err != nil {
		log.Println("Invalid join event:", err)
		return err
	}

	room := getOrCreateRoom(join.Name)
	client.room = room

	room.mu.Lock()
	room.clients[client] = true
	room.mu.Unlock()

	log.Println("Client joined room: ", join.Name)
	client.room.broadcast(event, client)
	return nil
}

func (m *Manager) servesWs(w http.ResponseWriter, r *http.Request) {
	fmt.Println("new connection recieved")
	roomName := r.URL.Query().Get("room")
	if roomName == "" {
		fmt.Println("Room name is required")
	}

	ws, err := upgrader.Upgrade(w, r, nil)

	if err != nil {
		fmt.Println(err)
		return
	}

	client := NewClient(ws, m)
	m.addClient(client)

	room := getOrCreateRoom(roomName)
	client.room = room

	room.mu.Lock()
	room.clients[client] = true
	room.mu.Unlock()

	go client.readMessage()
	go client.writeMessage()
}

func (m *Manager) addClient(client *Client) {
	fmt.Println("new client added")
	m.Lock()
	defer m.Unlock()
	m.clients[client] = true
}

func (m *Manager) removeClient(client *Client) {
	m.Lock()
	defer m.Unlock()
	if _, ok := m.clients[client]; ok {
		client.conn.Close()
		delete(m.clients, client)
		log.Println("Client removed from manager.")
	}
}

func (room *Room) removeClient(client *Client) {
	room.mu.Lock()
	defer room.mu.Unlock()
	if _, exists := room.clients[client]; exists {
		delete(room.clients, client)
		log.Println("Client removed from room:", room.name)
	}
}
