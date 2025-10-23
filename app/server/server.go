package server

import (
	"fmt"
	"log"
	"time"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	fileupload "github.com/chrollo-lucifer-12/excallidraw-backend/app/filleupload"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/ws"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type ServerOpts struct {
	Env          *dotenv.Env
	Database     *db.DB
	UploadClient *fileupload.UploadService
	Ws           *ws.RoomManager
}

type Server struct {
	env          *dotenv.Env
	db           *db.DB
	router       *gin.Engine
	uploadClient *fileupload.UploadService
	ws           *ws.RoomManager
}

func NewServer(opts ServerOpts) *Server {
	server := &Server{
		env:          opts.Env,
		db:           opts.Database,
		uploadClient: opts.UploadClient,
		ws:           opts.Ws,
	}
	server.router = gin.Default()
	return server
}

func (s *Server) Start() {

	s.router.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:     []string{"Origin", "Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
		MaxAge:           12 * time.Hour,
	}))

	s.RegisterRoutes(s.router)
	port := "8080"
	if s.env.PORT != "" {
		port = s.env.PORT
	}
	s.router.GET("/ws", func(c *gin.Context) {
		s.ws.HandleRequest(c.Writer, c.Request)
	})
	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running on %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
