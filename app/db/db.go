package db

import (
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
