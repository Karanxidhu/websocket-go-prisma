package main

import (
	"encoding/json"
	"log"

	"github.com/gorilla/websocket"
)

// var (
// 	pongWait   = 10 * time.Second
// 	pingPeriod = (pongWait * 9) / 10
// )

type ClientList map[*Client]bool

type Client struct {
	conn    *websocket.Conn
	manager *Manager
	room    *Room
	egress  chan Event
	name    string
}

func NewClient(conn *websocket.Conn, manager *Manager, userName string) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan Event, 10),
		name:    userName,
	}
}

func (client *Client) send(event Event) {
	select {
	case client.egress <- event:
		// Successfully sent to the egress channel
	default:
		log.Printf("Egress channel full, dropping message for client %s", client.name)
	}
}


func (client *Client) readMessage() {
	defer func() {
		log.Printf("Closing client connection for user: %s", client.name)
		client.manager.removeClient(client)
		if client.room != nil {
			client.room.removeClient(client)
		}
		client.conn.Close()
	}()

	client.conn.SetReadLimit(512)

	for {
		_, msg, err := client.conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Unexpected close error:", err)
			}
			log.Printf("Client %s disconnected, removing from manager and room.", client.name)
			break
		}

		var request Event
		if err := json.Unmarshal(msg, &request); err != nil {
			log.Println("Error decoding message:", err)
			continue
		}

		// Handle the event
		if err := client.manager.routeEvent(request, client); err != nil {
			log.Println("Error routing event:", err)
			continue
		}
	}

}

func (client *Client) writeMessage() {
	defer func() {
		log.Printf("Closing client connection for user: %s", client.name)
		client.manager.removeClient(client)
		if client.room != nil {
			client.room.removeClient(client)
		}
		client.conn.Close()
	}()
    for {
        select {
        case message, ok := <-client.egress:
            if !ok {
                log.Printf("Egress channel closed for client %s.", client.name)
                return
            }
    
            data, err := json.Marshal(message)
            if err != nil {
                log.Println("Error serializing message:", err)
                continue
            }
    
            if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
                log.Printf("Error writing message for client %s: %v", client.name, err)
                return
            }
        }
    }
    
}
