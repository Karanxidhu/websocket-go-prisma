package main

import (
	"encoding/json"
	"fmt"
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
}

func NewClient(conn *websocket.Conn, manager *Manager) *Client {
	return &Client{
		conn:    conn,
		manager: manager,
		egress:  make(chan Event),
	}
}

func (client *Client) readMessage() {
	fmt.Println("reading message")
	defer func() {
		client.manager.removeClient(client)
	}()

	// if err := client.conn.SetReadDeadline(time.Now().Add(pongWait)); err != nil {
	// 	log.Println("Error setting read deadline:", err)
	// }
	client.conn.SetReadLimit(512)

	// client.conn.SetPingHandler(client.pongHandler)

	for {
		_, msg, err := client.conn.ReadMessage()

		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Println("Error reading message:", err)
				client.room.removeClient(client)
			}
			break
		}
		fmt.Println("Received")
		var request Event

		if err := json.Unmarshal(msg, &request); err != nil {
			fmt.Println(err)
			break
		}
		if err := client.manager.routeEvent(request, client); err != nil {
			fmt.Println(err)
			break
		}
		client.room.broadcast(request, client)
	}

}

// func (client *Client) writeMessage() {
// 	fmt.Println("writing message")
// 	defer func() {
// 		client.manager.removeClient(client)
// 	}()

// 	// ticker := time.NewTicker(pingPeriod)

// 	for {
// 		select {
// 		case message, ok := <-client.egress:
// 			if !ok {
// 				if err := client.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
// 					log.Println(" connection closed: ", err)
// 				}
// 				return
// 			}
// 			data, err := json.Marshal(message)
// 			if err != nil {
// 				log.Println(err)
// 				return
// 			}
// 			if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
// 				log.Println("Error writing message:", err)
// 				return
// 			}

// 			// case <-ticker.C:
// 			// 	if err := client.conn.WriteMessage(websocket.PingMessage, []byte(``)); err != nil {
// 			// 		log.Println("Error writing message ping:", err)
// 			// 		return
// 			// 	}
// 		}
// 	}
// }
func (client *Client) writeMessage() {
	defer func() {
		client.manager.removeClient(client)
		if client.room != nil {
			client.room.removeClient(client)
		}
	}()

	for {
		select {
		case message, ok := <-client.egress:
			if !ok {
				// Connection has been closed
				if err := client.conn.WriteMessage(websocket.CloseMessage, nil); err != nil {
					log.Println("Connection closed:", err)
				}
				return
			}

			// Serialize and send the message to the WebSocket
			data, err := json.Marshal(message)
			if err != nil {
				log.Println("Error serializing message:", err)
				return
			}
			if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
				log.Println("Error writing message:", err)
				return
			}
		}
	}
}



// func (client *Client) pongHandler(appData string) error {
// 	client.conn.SetReadDeadline(time.Now().Add(pongWait))
// 	return nil
// }
