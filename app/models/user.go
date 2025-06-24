package models

import (
	"time"

	"github.com/go-playground/validator/v10"
)

type User struct {
	ID        uint `gorm:"primarykey"`
	CreatedAt time.Time
	UpdatedAt time.Time
	Username  string `json:"username" validate:"required,min=6,max=32" gorm:"unique;type:varchar(20)"`
	Password  string `json:"password,omitempty" validate:"required,min=6" gorm:"type:varchar(255)"`
	Fullname  string `json:"fullname" validate:"required,min=6" gorm:"type:varchar(100)"`
}

func (l User) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type UserSession struct {
	ID           uint `gorm:"primarykey"`
	CreatedAt    time.Time
	UpdatedAt    time.Time
	UserId       uint   `gorm:"not null;uniqueIndex:idx_user_session"`
	Token        string `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_session"`
	RefreshToken string `gorm:"type:varchar(255);not null;uniqueIndex:idx_user_session"`
}

func (l UserSession) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginRequest struct {
	Username string `json:"username" validate:"required"`
	Password string `json:"password" validate:"required"`
}

func (l LoginRequest) Validate() error {
	v := validator.New()
	return v.Struct(l)
}

type LoginResponse struct {
	Username          string `json:"username"`
	Fullname          string `json:"fullname"`
	Token             string `json:"token"`
	RefreshTokenToken string `json:"refresh_token"`
}
