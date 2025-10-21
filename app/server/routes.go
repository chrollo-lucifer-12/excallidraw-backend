package server

import "github.com/gin-gonic/gin"

func (s *Server) RegisterRoutes(r *gin.Engine) {
	api := r.Group("/api/v1")
	s.RegisterAuthRoutes(api)
	s.RegisterUserRoutes(api)
	s.RegisterWhiteboardRoutes(api)
}
