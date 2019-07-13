package model

const Tmpl = `package model

import (
	"cloud/lib/logger"
	"time"
)

{{ $md := .Model -}}
{{ $mds := .Models -}}

type {{$mds}} []{{$md}}

type {{$md}} struct {

}

type {{$md}}Unique struct {

}

func (md {{$md}}) TableName() string {
	return "{{.Table}}"
}

func New{{$md}}(version string, com uint64, bk, obj string, time time.Time, tag string) {{$md}} {

	u := {{$md}}Unique{Version: version, Com: com}
	return {{$md}}{
		{{$md}}Unique: u, 
		Bucket: bk, 
		Obj: obj, 
		Time: time, 
		Tag: tag,
	}
}

func GetAll{{$mds}}(tx *Tx) ({{$mds}}, error) {
	var mds {{$mds}}
	if err := tx.Find(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

{{range $x := .GetOneBy -}}

func Get{{$md}}By{{$x.Name}}(tx *Tx, {{$x.Params}}) ({{$md}}, error) {
	var md {{$md}}
	if _, err := tx.Where("{{$x.Query}}", {{$x.Args}}).Get(&md); err != nil {
		logger.Error(err)
		return md, err
	}
	return md, nil
}

{{end}}

{{range $x := .GetBy -}}

func Get{{$mds}}By{{$x.Name}}(tx *Tx, {{$x.Params}}) ({{$mds}}, error) {
	var mds {{$mds}}
	if _, err := tx.Where("{{$x.Query}}", {{$x.Args}}).Get(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

{{end}}

{{range $x := .GetByIn -}}

func Get{{$mds}}By{{$x.Name}}(tx *Tx, {{$x.Params}}) ({{$mds}}, error) {
	var mds {{$mds}}
	if _, err := tx.In("{{$x.Query}}", {{$x.Args}}).Get(&mds); err != nil {
		logger.Error(err)
		return mds, err
	}
	return mds, nil
}

{{end}}

func Count{{$md}}ByUnique(tx *Tx, c {{$md}}Unique) (int64, error) {

	n, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Count(&{{$md}}{})
	if err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func Insert{{$md}}(tx *Tx, mds ...*{{$md}}) error {
	for i := range mds {
		if _, err := tx.Insert(mds[i]); err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func Update{{$md}}(tx *Tx, md {{$md}}) error {
	if _, err := tx.Where("id = ?", md.ID).Update(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{$md}}(tx *Tx, md {{$md}}) error {

	if _, err := tx.Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{$md}}ById(tx *Tx, id uint64) error {
	var md {{$md}}
	if _, err := tx.Where("id = ?", id).Delete(&md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func Delete{{$md}}ByUnique(tx *Tx, c {{$md}}Unique) error {
	var md {{$md}}
	if _, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

`
