package models

import "time"

type User struct {
	ID            int
	Name          string
	Age           int
	Login         string
	Password      string
	CreatedAt     time.Time
	City          string
	TelegramLogin *string
	AvatarID      *int
	Role          string
}
