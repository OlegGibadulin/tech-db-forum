package models

type Forum struct {
	Title   string `json:"title" validate:"required"`
	User    string `json:"user" validate:"required,gte=3,lte=32"`
	Slug    string `json:"slug" validate:"required,gte=3,lte=64"`
	Posts   uint64 `json:"posts" validate:"eq=0"`
	Threads uint64 `json:"threads" validate:"eq=0"`
}
