package db

import (
	"crypto/sha1"
	"crypto/subtle"
	"encoding/hex"
	"errors"
	"log"
	"os"
	"strings"
	"time"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	"github.com/chrollo-lucifer-12/excallidraw-backend/app/util"
	"github.com/google/uuid"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const inactivityTimeoutSeconds = 60 * 60 * 24 * 10
const activityCheckIntervalSeconds = 60 * 60

type DBOpts struct {
	Env *dotenv.Env
}

type DB struct {
	Database *gorm.DB
}

func NewDB(opts DBOpts) (*DB, error) {
	databaseUrl := opts.Env.DATABASE_URL
	var err error
	db, err := gorm.Open(postgres.Open(databaseUrl), &gorm.Config{
		Logger: logger.New(
			log.New(os.Stdout, "\r\n", log.LstdFlags),
			logger.Config{
				SlowThreshold: time.Second,
				LogLevel:      logger.Info,
				Colorful:      true,
			},
		),
	})
	if err != nil {
		log.Fatalf("Failed to connect to database :%v\n", err)
		return nil, err
	}

	if err := db.Exec(`CREATE EXTENSION IF NOT EXISTS "uuid-ossp";`).Error; err != nil {
		log.Fatalf("failed to enable uuid-ossp: %v", err)
	}

	err = db.AutoMigrate(&User{}, &UserData{}, &Whiteboard{}, &WhiteboardElement{}, &Session{})
	if err != nil {
		log.Fatalf("Failed to migrate database :%v\n", err)
		return nil, err
	}

	return &DB{
		Database: db,
	}, nil
}

func (d *DB) FindUserByID(id string) (*User, error) {
	var user User
	if err := d.Database.Where("id = ?", id).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) FindUserByEmail(email string) (*User, error) {
	var user User
	if err := d.Database.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
}

func (d *DB) FindUserByUsername(username string) (*UserData, error) {
	userProfile := UserData{
		Username: username,
	}
	if err := d.Database.Where("username = ?", username).First(&userProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userProfile, nil
}

func (d *DB) CreateUser(email string, password string) error {
	user := User{
		Email:    email,
		Password: password,
	}
	if err := d.Database.Create(&user).Error; err != nil {
		return err
	}
	return nil
}

func (d *DB) GetUserProfile(user_id uuid.UUID) (*UserData, error) {
	userProfile := UserData{
		UserID: user_id,
	}

	if err := d.Database.Where("user_id = ?", user_id).First(&userProfile).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &userProfile, nil
}

func (d *DB) CreateUserProfile(birthDate time.Time, avatarUrl string, fullname, username string, user_id uuid.UUID) error {
	userProfile := UserData{
		UserID:    user_id,
		BirthDate: birthDate,
		AvatarUrl: avatarUrl,
		Fullname:  fullname,
		Username:  username,
	}
	if err := d.Database.Create(&userProfile).Error; err != nil {
		return err
	}
	return nil
}

func (d *DB) UpdateUserProfile(profile *UserData) error {
	if err := d.Database.Model(&UserData{}).
		Where("user_id = ?", profile.UserID).
		Updates(profile).Error; err != nil {
		return err
	}

	return nil
}

func (d *DB) CreateWhiteboard(admin_id uuid.UUID, name string, slug string) error {
	whiteboard := Whiteboard{
		Name:    name,
		AdminID: admin_id,
		Slug:    slug,
		Users: []User{
			{ID: admin_id},
		},
	}

	if err := d.Database.Create(&whiteboard).Error; err != nil {
		return err
	}
	return nil
}

func (d *DB) GetWhiteboardsByAdminID(admin_id uuid.UUID) ([]Whiteboard, error) {
	var whiteboards []Whiteboard
	if err := d.Database.Where("admin_id = ?", admin_id).Find(&whiteboards).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return whiteboards, nil
}

func (d *DB) CreateSession(user_id uuid.UUID) (string, error) {
	id, err := uuid.NewRandom()
	if err != nil {
		return "", err
	}
	secret, err := util.GenerateSecureRandomString()
	if err != nil {
		return "", err
	}

	token := id.String() + "." + secret

	secretHash := sha1.Sum([]byte(secret))

	session := Session{
		ID:             id,
		SecretHash:     hex.EncodeToString(secretHash[:]),
		LastVerifiedAt: time.Now(),
		UserId:         user_id,
	}

	if err := d.Database.Create(&session).Error; err != nil {
		return "", err
	}

	return token, nil
}

func (d *DB) DeleteSession(session_id uuid.UUID) {
	d.Database.Where("id = ?", session_id).Delete(&Session{})
}

func (d *DB) GetSession(session_id uuid.UUID) (*Session, error) {
	var session Session
	if err := d.Database.Where("id = ?", session_id).First(&session).Error; err != nil {
		return nil, err
	}

	if time.Since(session.LastVerifiedAt) >= time.Duration(inactivityTimeoutSeconds)*time.Second {
		d.DeleteSession(session_id)
		return nil, nil
	}
	return &session, nil
}

func (d *DB) UpdateSession(session *Session) {
	d.Database.Save(session)
}

func (d *DB) ValidateSessionToken(token string) (uuid.UUID, error) {
	tokenParts := strings.Split(token, ".")
	if len(tokenParts) != 2 {
		return uuid.Nil, errors.New("invalid token format")
	}
	sessionIdStr := tokenParts[0]
	sessionSecret := tokenParts[1]

	sessionId_UUID, err := util.ParseUUID(sessionIdStr)
	if err != nil {
		return uuid.Nil, err
	}

	session, err := d.GetSession(sessionId_UUID)
	if err != nil {
		return uuid.Nil, err
	}
	if session == nil {
		return uuid.Nil, errors.New("session expired or inactive")
	}
	tokenSecretHash := sha1.Sum([]byte(sessionSecret))
	storedHash, err := hex.DecodeString(session.SecretHash)
	if err != nil {
		return uuid.Nil, err
	}
	ok := subtle.ConstantTimeCompare(tokenSecretHash[:], storedHash)
	if ok != 1 {
		return uuid.Nil, errors.New("invalid session")
	}

	if time.Since(session.LastVerifiedAt) >= time.Duration(activityCheckIntervalSeconds)*time.Second {
		session.LastVerifiedAt = time.Now()
		d.UpdateSession(session)
	}

	return session.UserId, nil
}
