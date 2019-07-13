package model

import (
	"cloud/lib/logger"
	"time"
)

type MsFws []MsFw

type MsFw struct {
	ID         uint64 `json:"id" xorm:"'id' pk autoincr"`
	MsFwUnique `xorm:"extends"`
	Bucket     string    `json:"bucket"`
	Obj        string    `json:"obj"`
	Time       time.Time `json:"time"`
	Tag        string    `json:"tag"`
}

type MsFwUnique struct {
	Com     uint64 `json:"com"`
	Version string `json:"ver"`
}

func (md MsFw) TableName() string {
	return "cloud.msfw"
}

func NewMsFw(version string, com uint64, bk, obj string, time time.Time, tag string) MsFw {

	u := MsFwUnique{Version: version, Com: com}
	return MsFw{
		MsFwUnique: u,
		Bucket:     bk,
		Obj:        obj,
		Time:       time,
		Tag:        tag,
	}
}

func GetAllMsFws(tx *Tx) (MsFws, error) {
	var mds MsFws
	if err := tx.Find(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func GetMsFwById(tx *Tx, id uint64) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("id = ?", id).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func GetMsFwByUnique(tx *Tx, c MsFwUnique) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func GetMsFwByCom(tx *Tx, com uint64) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("com = ?", com).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func GetMsFwsByComs(tx *Tx, coms []uint64) (MsFws, error) {
	var mds MsFws
	if _, err := tx.In("com = ?", coms).Get(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func CountMsFwByUnique(tx *Tx, c MsFwUnique) (int64, error) {

	n, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Count(&MsFw{})
	if err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func InsertMsFw(tx *Tx, mds ...*MsFw) error {
	for i := range mds {
		if _, err := tx.Insert(mds[i]); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func UpdateMsFw(tx *Tx, md MsFw) error {
	if _, err := tx.Where("id = ?", md.ID).Update(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func DeleteMsFw(tx *Tx, md MsFw) error {

	if _, err := tx.Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func DeleteMsFwById(tx *Tx, id uint64) error {
	var md MsFw
	if _, err := tx.Where("id = ?", id).Delete(&md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func DeleteMsFwByUnique(tx *Tx, c MsFwUnique) error {
	var md MsFw
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

// type MsFwStore interface {
// 	InsertMsFw(model.Tx, MsFw) (MsFw, error)
// 	MsFws(model.Tx) (MsFws, error)
// 	MsFwByID(model.Tx, uint64) (MsFw, error)
// 	MsFwByUnique(model.Tx, MsFwUnique) (MsFw, error)
// 	MsFwsByComs(model.Tx, []uint64) (MsFws, error)
// 	MsFwsByCom(model.Tx, uint64) (MsFws, error)
// 	MsFwsByComAndShare(model.Tx, uint64) (MsFws, error)
// 	MsFwsByGroup(model.Tx, uint64) (MsFws, error)
// 	DeleteMsFw(model.Tx, uint64) error
// 	DeleteMsFwByUnique(model.Tx, MsFwUnique) error
// 	MsFwCompanyByUnique(model.Tx, MsFwUnique) (string, error)
// 	CountMsFwByUnique(model.Tx, MsFwUnique) (uint64, error)
// }
