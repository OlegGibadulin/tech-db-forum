package models

type User struct {
	Nickname string `json:"nickname" validate:"gte=3,lte=32"`
	Fullname string `json:"fullname" validate:"required,gte=3,lte=32"`
	Email    string `json:"email" validate:"required,email,lte=64"`
	About    string `json:"about"`
}
