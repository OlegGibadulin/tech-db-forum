package models

import (
	"time"
)

type Filter struct {
	Limit uint64    `query:"limit" validate:"gte=0"`
	Since time.Time `query:"since"`
	Desc  bool      `query:"desc"`
}
