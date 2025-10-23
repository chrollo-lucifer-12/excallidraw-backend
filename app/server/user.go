package server

import (
	"net/http"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterUserRoutes(r *gin.RouterGroup) {
	user := r.Group("/user", s.userMiddleware)
	{
		user.POST("/create-profile", s.createProfileHandler)
		user.POST("/update-profile", s.updateProfileHandler)
		user.GET("/me", s.getUserProfileHandler)
		user.POST("/upload-avatar", s.uploadAvatarHandler)
		user.GET("/whiteboards", s.getWhiteboardsHandler)
	}
}

type CreateProfileRequest struct {
	BirthDate string `json:"birthdate"`
	AvatarUrl string `json:"avatarUrl"`
	Fullname  string `json:"fullname" binding:"required"`
	Username  string `json:"username" binding:"required"`
}

type UpdateProfileRequest struct {
	BirthDate string `json:"birthdate"`
	AvatarUrl string `json:"avatarUrl"`
	Fullname  string `json:"fullname"`
	Username  string `json:"username"`
}

func (s *Server) createProfileHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	req, err := util.BindJSON[CreateProfileRequest](c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	findUser, err := s.db.FindUserByUsername(req.Username)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if findUser != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Username already exists"})
		return
	}

	user_id_uuid, err := util.ParseUUID(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	profile, err := s.db.GetUserProfile(user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if profile != nil {
		c.JSON(http.StatusConflict, gin.H{"error": "Profile already exists"})
		//	panic(err)
	}

	birthdate_time, err := util.ParseTime(req.BirthDate)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	err = s.db.CreateUserProfile(birthdate_time, req.AvatarUrl, req.Fullname, req.Username, user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "Profile created successfully"})
}

func (s *Server) updateProfileHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	req, err := util.BindJSON[UpdateProfileRequest](c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	user_id_uuid, err := util.ParseUUID(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	profile, err := s.db.GetUserProfile(user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if profile == nil {
		c.JSON(http.StatusNoContent, gin.H{"error": "Profile doesn't exist"})
		//	panic(err)
	}

	if req.BirthDate != "" {
		birthdateTime, err := util.ParseTime(req.BirthDate)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid birth date format"})
			return
		}
		profile.BirthDate = birthdateTime
	}

	if req.AvatarUrl != "" {
		profile.AvatarUrl = req.AvatarUrl
	}

	if req.Fullname != "" {
		profile.Fullname = req.Fullname
	}

	if req.Username != "" {
		profile.Username = req.Username
	}

	// Save updated profile
	if err := s.db.UpdateUserProfile(profile); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Profile updated successfully"})
}

func (s *Server) getUserProfileHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	user_id_uuid, err := util.ParseUUID(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}

	userProfile, err := s.db.GetUserProfile(user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	if userProfile == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Profile not found"})
		//	panic(err)
	}

	c.JSON(http.StatusFound, gin.H{"message": "Profile found", "user": userProfile})
}

func (s *Server) uploadAvatarHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "No file is received"})
		return
	}
	src, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Cannot open file"})
		panic(err)
	}
	defer src.Close()
	err = s.uploadClient.UploadFile("avatars", file.Filename, file.Size, src)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "error uploading file"})
		panic(err)
	}

	c.JSON(http.StatusCreated, gin.H{"message": "File uplaoded"})
}

func (s *Server) getWhiteboardsHandler(c *gin.Context) {
	user_id, ok := c.Get("user_id")
	if !ok {
		return
	}

	user_id_uuid, err := util.ParseUUID(user_id.(string))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	whiteboards, err := s.db.GetWhiteboardsByAdminID(user_id_uuid)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		panic(err)
	}
	c.JSON(http.StatusOK, gin.H{"whiteboards": whiteboards})
}
