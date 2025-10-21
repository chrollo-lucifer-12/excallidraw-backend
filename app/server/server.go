package server

import (
	"fmt"
	"log"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	fileupload "github.com/chrollo-lucifer-12/excallidraw-backend/app/filleupload"
	"github.com/gin-gonic/gin"
)

type ServerOpts struct {
	Env          *dotenv.Env
	Database     *db.DB
	UploadClient *fileupload.UploadService
}

type Server struct {
	env          *dotenv.Env
	db           *db.DB
	router       *gin.Engine
	uploadClient *fileupload.UploadService
}

func NewServer(opts ServerOpts) *Server {
	server := &Server{
		env:          opts.Env,
		db:           opts.Database,
		uploadClient: opts.UploadClient,
	}
	server.router = gin.Default()
	return server
}

func (s *Server) Start() {
	s.RegisterRoutes(s.router)
	port := "8080"
	if s.env.PORT != "" {
		port = s.env.PORT
	}

	addr := fmt.Sprintf(":%s", port)
	log.Printf("🚀 Server running on %s", addr)
	if err := s.router.Run(addr); err != nil {
		log.Fatalf("❌ Failed to start server: %v", err)
	}
}
