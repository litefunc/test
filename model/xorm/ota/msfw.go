package model

import (
	"cloud/lib/logger"
	"cloud/server/ota/service/v2/msfw"
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

func (md MsFwUnique) ServiceData() msfw.MsFwUnique {

	return msfw.MsFwUnique{
		Version: md.Version,
		Com:     md.Com,
	}
}

func (md MsFw) ServiceData() msfw.MsFw {

	return msfw.MsFw{
		MsFwUnique: md.MsFwUnique.ServiceData(),
		Bucket:     md.Bucket,
		Obj:        md.Obj,
		Time:       md.Time,
		Tag:        md.Tag,
	}
}

func (md MsFws) ServiceData() msfw.MsFws {

	var mds msfw.MsFws
	for _, v := range md {
		mds = append(mds, v.ServiceData())
	}
	return mds
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

func (md *MsFw) SetData(data msfw.MsFw) {

	u := MsFwUnique{Version: data.Version, Com: data.Com}
	md = &MsFw{
		MsFwUnique: u,
		Bucket:     data.Bucket,
		Obj:        data.Obj,
		Time:       data.Time,
		Tag:        data.Tag,
	}
}

func (md *MsFwUnique) SetData(data msfw.MsFwUnique) {
	md = &MsFwUnique{Version: data.Version, Com: data.Com}
}

func (pg Pg) GetAllMsFws(tx *Tx) (msfw.MsFws, error) {
	mds, err := GetAllMsFws(tx)
	return mds.ServiceData(), err
}

func GetAllMsFws(tx *Tx) (MsFws, error) {
	var mds MsFws
	if err := tx.Find(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func (pg Pg) GetMsFwById(tx *Tx, id uint64) (msfw.MsFw, error) {
	md, err := GetMsFwById(tx, id)
	return md.ServiceData(), err
}

func GetMsFwById(tx *Tx, id uint64) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("id = ?", id).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func (pg Pg) GetMsFwByUnique(tx *Tx, u msfw.MsFwUnique) (msfw.MsFw, error) {
	var by *MsFwUnique
	by.SetData(u)
	md, err := GetMsFwByUnique(tx, *by)
	return md.ServiceData(), err
}

func GetMsFwByUnique(tx *Tx, u MsFwUnique) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("com = ? AND version = ?", u.Com, u.Version).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func (pg Pg) GetMsFwByCom(tx *Tx, com uint64) (msfw.MsFw, error) {
	md, err := GetMsFwByCom(tx, com)
	return md.ServiceData(), err
}

func GetMsFwByCom(tx *Tx, com uint64) (MsFw, error) {
	var md MsFw
	if _, err := tx.Where("com = ?", com).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func (pg Pg) GetMsFwsByComs(tx *Tx, coms []uint64) (msfw.MsFws, error) {
	mds, err := GetMsFwsByComs(tx, coms)
	return mds.ServiceData(), err
}

func GetMsFwsByComs(tx *Tx, coms []uint64) (MsFws, error) {
	var mds MsFws
	if _, err := tx.In("com = ?", coms).Get(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func (pg Pg) CountMsFwByUnique(tx *Tx, u msfw.MsFwUnique) (uint64, error) {
	var by *MsFwUnique
	by.SetData(u)
	i, err := CountMsFwByUnique(tx, *by)
	return uint64(i), err
}

func CountMsFwByUnique(tx *Tx, c MsFwUnique) (int64, error) {

	n, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Count(&MsFw{})
	if err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func (pg Pg) InsertMsFw(tx *Tx, mds ...*msfw.MsFw) error {
	var bys []*MsFw
	for i := range mds {
		var by *MsFw
		by.SetData(*mds[i])
		bys = append(bys, by)
	}
	return InsertMsFw(tx, bys...)
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

func (pg Pg) UpdateMsFw(tx *Tx, md msfw.MsFw) error {
	var by *MsFw
	by.SetData(md)
	return UpdateMsFw(tx, *by)
}

func UpdateMsFw(tx *Tx, md MsFw) error {
	if _, err := tx.Where("id = ?", md.ID).Update(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) DeleteMsFw(tx *Tx, md msfw.MsFw) error {
	var by *MsFw
	by.SetData(md)
	return DeleteMsFw(tx, *by)
}

func DeleteMsFw(tx *Tx, md MsFw) error {

	if _, err := tx.Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) DeleteMsFwById(tx *Tx, id uint64) error {
	return DeleteMsFwById(tx, id)
}

func DeleteMsFwById(tx *Tx, id uint64) error {
	var md MsFw
	if _, err := tx.Where("id = ?", id).Delete(&md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) DeleteMsFwByUnique(tx *Tx, md msfw.MsFwUnique) error {
	var by *MsFwUnique
	by.SetData(md)
	return DeleteMsFwByUnique(tx, *by)
}

func DeleteMsFwByUnique(tx *Tx, c MsFwUnique) error {
	var md MsFw
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}
