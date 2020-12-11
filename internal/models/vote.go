package models

type Vote struct {
	Nickname string `json:"nickname" validate:"required,gte=3,lte=32"`
	Voice    int    `json:"voice"`
}
