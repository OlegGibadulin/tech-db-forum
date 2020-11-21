package models

import (
	"time"
)

type Post struct {
	ID       uint64    `json:"id" validate:"isdefault"`
	Parent   uint64    `json:"parent" validate:"gte=0"`
	Author   string    `json:"author" validate:"required,gte=3,lte=32"`
	Message  string    `json:"message" validate:"required"`
	IsEdited bool      `json:"isEdited"`
	Forum    string    `json:"forum" validate:"gte=3,lte=64"`
	Thread   uint64    `json:"thread" validate:"gte=0"`
	Created  time.Time `json:"created"`
}
