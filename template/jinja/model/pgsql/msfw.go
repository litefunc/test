package model

import (
	"cloud/lib/logger"
	"test/template/jinja/out/service/msfw"
	"time"
)

type Msfw struct {
	ID uint64 `json:"id"`
	MsfwUnique
	Bucket string    `json:"bucket"`
	Obj    string    `json:"obj"`
	Time   time.Time `json:"time"`
	Tag    string    `json:"tag"`
}

type MsfwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver"`
}

type Msfws []Msfw

func NewMsfw(version string, com uint64, bk, obj string, time time.Time, tag string) Msfw {
	msfwkey := MsfwUnique{Version: version, Com: com}
	return Msfw{MsfwUnique: msfwkey, Bucket: bk, Obj: obj, Time: time, Tag: tag}
}

func NewMsfwUnique(com uint64, version string) MsfwUnique {
	return MsfwUnique{Com: com, Version: version}
}

func (md Msfw) From(m msfw.Msfw) Msfw {
	u := MsfwUnique{m.MsfwUnique.Com, m.MsfwUnique.Version}
	return Msfw{m.ID, u, m.Bucket, m.Obj, m.Time, m.Tag}
}

func (md Msfw) Msfw() msfw.Msfw {
	u := msfw.MsfwUnique{md.MsfwUnique.Com, md.MsfwUnique.Version}
	return msfw.Msfw{md.ID, u, md.Bucket, md.Obj, md.Time, md.Tag}
}

func (md Msfws) Msfws() msfw.Msfws {
	mds := make(msfw.Msfws, len(md), len(md))
	for i := range md {
		mds[i] = md[i].Msfw()
	}
	return mds
}

func (md MsfwUnique) From(m msfw.MsfwUnique) MsfwUnique {
	return MsfwUnique{m.Com, m.Version}
}

func (md MsfwUnique) MsfwUnique() msfw.MsfwUnique {
	return msfw.MsfwUnique{md.Com, md.Version}
}

func (db DB) InsertMsfw(in msfw.Msfw) (uint64, error) {

	md := Msfw{}.From(in)
	if err := db.Insert(md).Returning("id").Run().Scan(&md.ID); err != nil {
		logger.Error(err)
		return 0, err
	}
	return md.ID, nil
}
func (db DB) GetMsfws(condition map[string]interface{}, cols ...string) (msfw.Msfws, error) {
	var md Msfws
	c, v := arg(condition)
	if err := db.Select(&md, cols...).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return nil, err
	}
	return md.Msfws(), nil
}
func (db DB) GetOneMsfw(condition map[string]interface{}, cols ...string) (msfw.Msfw, error) {
	var md Msfw
	c, v := arg(condition)
	if err := db.Select(&md, cols...).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return msfw.Msfw{}, err
	}
	return md.Msfw(), nil
}
func (db DB) UpdateMsfw(set, condition map[string]interface{}) error {
	var md Msfw
	c, v := arg(condition)
	if err := db.Update(&md).Set(c, v...).Run(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (db DB) DeleteMsfw(condition map[string]interface{}) error {
	var md Msfw
	c, v := arg(condition)
	if err := db.Delete(&md).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (db DB) CountMsfw(condition map[string]interface{}) (uint64, error) {
	var md Msfw
	c, v := arg(condition)
	var n uint64
	if err := db.Count(&md).Where(c, v...).Run().Scan(&n); err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func (tx Tx) InsertMsfw(in msfw.Msfw) (uint64, error) {

	md := Msfw{}.From(in)
	if err := tx.Insert(md).Returning("id").Run().Scan(&md.ID); err != nil {
		logger.Error(err)
		return 0, err
	}
	return md.ID, nil
}
func (tx Tx) GetMsfws(condition map[string]interface{}, cols ...string) (msfw.Msfws, error) {
	var md Msfws
	c, v := arg(condition)
	if err := tx.Select(&md, cols...).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return nil, err
	}
	return md.Msfws(), nil
}
func (tx Tx) GetOneMsfw(condition map[string]interface{}, cols ...string) (msfw.Msfw, error) {
	var md Msfw
	c, v := arg(condition)
	if err := tx.Select(&md, cols...).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return msfw.Msfw{}, err
	}
	return md.Msfw(), nil
}
func (tx Tx) UpdateMsfw(set, condition map[string]interface{}) error {
	var md Msfw
	c, v := arg(condition)
	if err := tx.Update(&md).Set(c, v...).Run(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (tx Tx) DeleteMsfw(condition map[string]interface{}) error {
	var md Msfw
	c, v := arg(condition)
	if err := tx.Delete(&md).Where(c, v...).Run(); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
func (tx Tx) CountMsfw(condition map[string]interface{}) (uint64, error) {
	var md Msfw
	c, v := arg(condition)
	var n uint64
	if err := tx.Count(&md).Where(c, v...).Run().Scan(&n); err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}
