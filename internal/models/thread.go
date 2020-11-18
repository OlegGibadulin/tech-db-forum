package models

import (
	"time"
)

type Thread struct {
	ID      uint64    `json:"id" validate:"isdefault"`
	Title   string    `json:"title" validate:"required,gte=3,lte=64"`
	Author  string    `json:"author" validate:"required,gte=3,lte=32"`
	Forum   string    `json:"forum" validate:"gte=3,lte=64"`
	Message string    `json:"message" validate:"required"`
	Votes   uint64    `json:"votes" validate:"gte=0"`
	Slug    string    `json:"slug" validate:"omitempty,gte=3,lte=64"`
	Created time.Time `json:"created"`
}
