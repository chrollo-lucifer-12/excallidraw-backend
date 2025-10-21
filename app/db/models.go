package db

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID          uuid.UUID    `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Email       string       `gorm:"uniqueIndex;size:100;not null"`
	Password    string       `gorm:"not null"`
	CreatedAt   time.Time    `gorm:"autoCreateTime"`
	UpdatedAt   time.Time    `gorm:"autoUpdateTime"`
	Whiteboards []Whiteboard `gorm:"many2many:user_whiteboards;"`
}

type UserData struct {
	ID        uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	UserID    uuid.UUID `gorm:"type:uuid;not null;uniqueIndex"`
	BirthDate time.Time
	AvatarUrl string
	Fullname  string    `gorm:"not null"`
	Username  string    `gorm:"uniqueIndex;size:20;not null"`
	CreatedAt time.Time `gorm:"autoCreateTime"`
	UpdatedAt time.Time `gorm:"autoUpdateTime"`
}

type Whiteboard struct {
	ID      uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	Name    string    `gorm:"size:50;not null"`
	Slug    string    `gorm:"size:10;uniqueIndex"`
	AdminID uuid.UUID `gorm:"type:uuid;not null"`

	Users    []User `gorm:"many2many:user_whiteboards;"`
	Admin    User   `gorm:"foreignKey:AdminID"`
	Elements []WhiteboardElement
}

type WhiteboardElement struct {
	ID           uuid.UUID `gorm:"type:uuid;default:uuid_generate_v4();primaryKey"`
	WhiteboardID uuid.UUID `gorm:"type:uuid;not null;index"`
	Type         string
	Data         string
	CreatedAt    time.Time `gorm:"autoCreateTime"`
	UpdatedAt    time.Time `gorm:"autoUpdateTime"`

	Whiteboard Whiteboard `gorm:"foreignKey:WhiteboardID"`
}
