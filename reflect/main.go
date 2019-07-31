package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"reflect"
	"time"
)

type MsFws []MsFw

type MsFw struct {
	TableName struct{} `json:"-" sql:"cloud.msfw"`
	ID        uint64   `json:"id" sql:",pk"`
	MsFwUnique
	Bucket string      `json:"bucket"`
	Obj    string      `json:"obj"`
	Time   time.Time   `json:"time"`
	Tag    null.String `json:"tag"`
}

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver" sql:",notnull"`
}

func NewMsFw(version string, com uint64, bk, obj string, time time.Time, tag string) MsFw {

	u := MsFwUnique{Version: version, Com: com}
	return MsFw{
		MsFwUnique: u,
		Bucket:     bk,
		Obj:        obj,
		Time:       time,
		Tag:        null.NewString(tag),
	}
}

func main() {
	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	t := reflect.TypeOf(md1)
	logger.Debug(t.Kind())
	logger.Debug(t.String())

	field1 := t.Field(1)
	tag1 := field1.Tag.Get("sql")
	logger.Debug(tag1)

	for i := 0; i < t.NumField(); i++ {
		logger.Debug(t.Field(i))
		logger.Debug(t.Field(i).Type)
	}

}
