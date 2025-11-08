package ws

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

type Message struct {
	Type    string `json:"type"`
	Payload any    `json:"payload"`
}

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsOpts struct {
	database *db.DB
}

type Client struct {
	conn   *websocket.Conn
	userId string
}

type WebSocketManager struct {
	database *db.DB
	rooms    map[string][]*Client
}

func NewWs(database *db.DB) *WebSocketManager {
	return &WebSocketManager{
		database: database,
		rooms:    make(map[string][]*Client),
	}
}

func (ws *WebSocketManager) WsHandler(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		fmt.Println("Error upgrading:", err)
		return
	}
	go ws.handleConnection(conn)
}

func (ws *WebSocketManager) handleConnection(conn *websocket.Conn) {
	defer conn.Close()
	var client *Client
	for {
		_, msg, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		var message Message
		if err := json.Unmarshal(msg, &message); err != nil {
			fmt.Println("JSON parse error:", err)
			continue
		}

		switch message.Type {
		case "join":
			data, _ := json.Marshal(message.Payload)
			var joinData struct {
				RoomId string `json:"roomId"`
				UserId string `json:"userId"`
			}
			if err := json.Unmarshal(data, &joinData); err != nil {
				fmt.Println("Invalid join payload:", err)
				continue
			}

			client = &Client{
				conn:   conn,
				userId: joinData.UserId,
			}

			ws.rooms[joinData.RoomId] = append(ws.rooms[joinData.RoomId], client)
			ws.broadcastToRoom(joinData.RoomId, joinData.UserId, Message{
				Type:    "user_joined",
				Payload: joinData.UserId,
			})
			fmt.Printf("User %s joined room %s\n", joinData.UserId, joinData.RoomId)
		case "shapes":
			data, _ := json.Marshal(message.Payload)
			var shapesData struct {
				RoomId string `json:"roomId"`
				UserId string `json:"userId"`
				Shapes string `json:"shapes"`
			}
			if err := json.Unmarshal(data, &shapesData); err != nil {
				fmt.Println("Invalid join payload:", err)
				continue
			}

			ws.broadcastToRoom(shapesData.RoomId, shapesData.UserId, Message{
				Type: "shapes_update",
				Payload: map[string]any{
					"from":   shapesData.UserId,
					"shapes": shapesData.Shapes,
				},
			})
		}
	}
}

func (ws *WebSocketManager) broadcastToRoom(roomId string, userId string, message Message) {
	clients, ok := ws.rooms[roomId]
	if !ok {
		return
	}

	data, err := json.Marshal(message)
	if err != nil {
		fmt.Println("Error marshalling message:", err)
		return
	}

	for _, client := range clients {
		if client.userId == userId {
			continue
		}
		if err := client.conn.WriteMessage(websocket.TextMessage, data); err != nil {
			fmt.Println("Error writing to client:", err)
		}
	}
}
