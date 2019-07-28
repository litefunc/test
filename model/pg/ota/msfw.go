package model

import (
	"cloud/lib/logger"

	"encoding/json"
	"cloud/lib/null"
	"time"
)

type MsFws []MsFw

type MsFw struct {
	TableName struct{} `json:"-" sql:"cloud.msfw"`
	ID        uint64   `json:"id" sql:",pk"`
	MsFwUnique
	Bucket string         `json:"bucket"`
	Obj    string         `json:"obj"`
	Time   time.Time      `json:"time"`
	Tag    null.String `json:"tag"`
}

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver"`
}

func GetAllMsFws(tx *Tx) (MsFws, error) {
	var mds MsFws
	if err := tx.Select(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
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

func (md MsFw) Json() []byte {
	by, _ := json.Marshal(md)
	return by
}

func (md MsFws) Json() []byte {
	by, _ := json.Marshal(md)
	return by
}
