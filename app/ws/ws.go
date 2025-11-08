package ws

import (
	"fmt"
	"net/http"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

type WsOpts struct {
	database *db.DB
}

type WebSocketManager struct {
	database *db.DB
}

func NewWs(database *db.DB) *WebSocketManager {
	return &WebSocketManager{
		database: database,
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
	for {
		_, message, err := conn.ReadMessage()
		if err != nil {
			fmt.Println("Error reading message:", err)
			break
		}

		fmt.Printf("Received: %s\n", message)

		if err := conn.WriteMessage(websocket.TextMessage, message); err != nil {
			fmt.Println("Error writing message:", err)
			break
		}
	}
}
