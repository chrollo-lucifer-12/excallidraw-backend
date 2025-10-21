package server

import (
	"fmt"
	"log"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/db"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	"github.com/gin-gonic/gin"
)

type ServerOpts struct {
	Env      *dotenv.Env
	Database *db.DB
}

type Server struct {
	env    *dotenv.Env
	db     *db.DB
	router *gin.Engine
}

func NewServer(opts ServerOpts) *Server {
	server := &Server{
		env: opts.Env,
		db:  opts.Database,
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
