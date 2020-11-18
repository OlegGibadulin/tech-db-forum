package models

type Filter struct {
	Limit uint64 `query:"limit" validate:"gte=0"`
	Desc  bool   `query:"desc"`
}
