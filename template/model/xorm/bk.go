package model

const bk = `package model

import (
	"cloud/lib/logger"
	"time"
)

type {{.Models}} []{{.Model}}

type {{.Model}} struct {

}

type {{.Model}}Unique struct {

}

func (md {{.Model}}) TableName() string {
	return "{{.Table}}"
}

func New{{.Model}}(version string, com uint64, bk, obj string, time time.Time, tag string) {{.Model}} {

	u := {{.Model}}Unique{Version: version, Com: com}
	return {{.Model}}{
		{{.Model}}Unique: u, 
		Bucket: bk, 
		Obj: obj, 
		Time: time, 
		Tag: tag,
	}
}

func GetAll{{.Models}}(tx *Tx) ({{.Models}}, error) {
	var mds {{.Models}}
	if err := tx.Find(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func Get{{.Model}}ById(tx *Tx, id uint64) ({{.Model}}, error) {
	var md {{.Model}}
	if _, err := tx.Where("id = ?", id).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func Get{{.Model}}ByUnique(tx *Tx, c {{.Model}}Unique) ({{.Model}}, error) {
	var md {{.Model}}
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func Get{{.Model}}ByCom(tx *Tx, com uint64) ({{.Model}}, error) {
	var md {{.Model}}
	if _, err := tx.Where("com = ?", com).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

func Get{{.Models}}ByComs(tx *Tx, coms []uint64) ({{.Models}}, error) {
	var mds {{.Models}}
	if _, err := tx.In("com = ?", coms).Get(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

func Count{{.Model}}ByUnique(tx *Tx, c {{.Model}}Unique) (int64, error) {

	n, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Count(&{{.Model}}{})
	if err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func Insert{{.Model}}(tx *Tx, mds ...*{{.Model}}) error {
	for i := range mds {
		if _, err := tx.Insert(mds[i]); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func Update{{.Model}}(tx *Tx, md {{.Model}}) error {
	if _, err := tx.Where("id = ?", md.ID).Update(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{.Model}}(tx *Tx, md {{.Model}}) error {

	if _, err := tx.Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{.Model}}ById(tx *Tx, id uint64) error {
	var md {{.Model}}
	if _, err := tx.Where("id = ?", id).Delete(&md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{.Model}}ByUnique(tx *Tx, c {{.Model}}Unique) error {
	var md {{.Model}}
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

`
