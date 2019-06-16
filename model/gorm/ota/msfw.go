package model

import (
	"time"

	_ "github.com/jinzhu/gorm/dialects/postgres"
)

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver"`
}

type MsFw struct {
	ID uint64 `json:"id"`
	MsFwUnique
	Bucket string    `json:"bucket"`
	Obj    string    `json:"obj"`
	Time   time.Time `json:"time"`
	Tag    string    `json:"tag"`
}

func (MsFw) TableName() string {
	return "msfw"
}
