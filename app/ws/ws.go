package ws

import (
	"encoding/json"
	"fmt"
	"net/http"
	"sync"

	"github.com/olahol/melody"
)

type Message struct {
	UserID string `json:"userId"`
	Text   string `json:"text"`
}

type RoomManager struct {
	m     *melody.Melody
	rooms map[string]map[*melody.Session]string
	mu    sync.Mutex
}

func NewRoomManager() *RoomManager {
	rm := &RoomManager{
		m:     melody.New(),
		rooms: make(map[string]map[*melody.Session]string),
	}

	rm.m.HandleConnect(func(s *melody.Session) {
		if s.Request == nil {
			s.CloseWithMsg([]byte("no request info"))
			return
		}

		q := s.Request.URL.Query()
		roomID := q.Get("roomId")
		userID := q.Get("userId")

		if roomID == "" || userID == "" {
			s.CloseWithMsg([]byte("missing roomId or userId"))
			return
		}

		rm.AddUser(roomID, userID, s)
		fmt.Printf("User %s joined room %s\n", userID, roomID)
	})

	rm.m.HandleDisconnect(func(s *melody.Session) {
		rm.RemoveUser(s)
	})

	rm.m.HandleMessage(func(s *melody.Session, msg []byte) {
		if s.Request == nil {
			return
		}

		q := s.Request.URL.Query()
		roomID := q.Get("roomId")
		userID := q.Get("userId")

		if roomID == "" || userID == "" {
			return
		}

		out, _ := json.Marshal(Message{
			UserID: userID,
			Text:   string(msg),
		})

		rm.BroadcastToRoom(roomID, out)
	})

	return rm
}

func (rm *RoomManager) AddUser(roomID, userID string, s *melody.Session) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	if rm.rooms[roomID] == nil {
		rm.rooms[roomID] = make(map[*melody.Session]string)
	}
	rm.rooms[roomID][s] = userID
}

func (rm *RoomManager) RemoveUser(s *melody.Session) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for roomID, sessions := range rm.rooms {
		if _, ok := sessions[s]; ok {
			delete(sessions, s)
			if len(sessions) == 0 {
				delete(rm.rooms, roomID)
			}
			break
		}
	}
}

func (rm *RoomManager) BroadcastToRoom(roomID string, msg []byte) {
	rm.mu.Lock()
	defer rm.mu.Unlock()

	for s := range rm.rooms[roomID] {
		s.Write(msg)
	}
}

func (rm *RoomManager) HandleRequest(w http.ResponseWriter, r *http.Request) {
	if rm.m == nil {
		fmt.Println("Error: Melody instance is nil")
		return
	}
	rm.m.HandleRequest(w, r)
}
