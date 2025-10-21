package db

import (
	"errors"
	"log"
	"os"
	"time"

	"github.com/chrollo-lucifer-12/excallidraw-backend/app/dotenv"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

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

	err = db.AutoMigrate(&User{}, &UserData{}, &Whiteboard{}, &WhiteboardElement{})
	if err != nil {
		log.Fatalf("Failed to migrate database :%v\n", err)
		return nil, err
	}

	return &DB{
		Database: db,
	}, nil
}

func (d *DB) FindUser(email string) (*User, error) {
	var user User
	if err := d.Database.Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return nil, nil
		}
		return nil, err
	}
	return &user, nil
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
