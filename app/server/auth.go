package server

import (
	"net/http"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/gin-gonic/gin"
)

type SingupRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

type LoginRequest struct {
	Email    string `json:"email" binding:"required,email"`
	Password string `json:"password" binding:"required,min=6"`
}

func (s *Server) RegisterAuthRoutes(r *gin.RouterGroup) {
	auth := r.Group("/auth")
	{
		auth.POST("/login", s.loginHandler)
		auth.POST("/singup", s.singupHandler)
	}
}

func (s *Server) loginHandler(c *gin.Context) {
	req, err := util.BindJSON[LoginRequest](c)
	if err != nil {

		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password missing"})
		panic(err)
	}
	email := req.Email
	password := req.Password

	findUser, err := s.db.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		panic(err)
	}
	if findUser == nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User doesn't exist"})
		return
	}

	ok := util.CheckPassword(findUser.Password, password)

	if !ok {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Wrong password"})
		return
	}

	token, err := util.CreateToken(findUser.ID.String(), "sahil")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "login successful", "token": token})
}

func (s *Server) singupHandler(c *gin.Context) {
	req, err := util.BindJSON[SingupRequest](c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Email or password missing"})
		return
	}
	email := req.Email
	password := req.Password

	findUser, err := s.db.FindUserByEmail(email)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}
	if findUser != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
		return
	}

	hashPassword, err := util.HashPassword(password)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	err = s.db.CreateUser(email, hashPassword)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating user"})
		return
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Created user successfully"})
}
