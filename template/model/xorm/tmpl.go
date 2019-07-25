package model

const Tmpl = `package model

import (
	"cloud/lib/logger"
	"time"
)

{{ $md := .Model -}}
{{ $mds := .Models -}}
{{ $pkg := .ServicePackage -}}


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
	}
}

func (md {{$md}}Unique) ServiceData() {{$pkg}}.{{$md}}Unique {

	return {{$pkg}}.{{$md}}Unique{
	}
}

func (md {{$md}}) ServiceData() {{$pkg}}.{{$md}} {

	return {{$pkg}}.{{$md}}{
		{{$md}}Unique: md.{{$md}}Unique.ServiceData(),
	}
}

func (md {{$md}}s) ServiceData() {{$pkg}}.{{$md}}s {

	var mds {{$pkg}}.{{$md}}s
	for _, v := range md {
		mds = append(mds, v.ServiceData())
	}
	return mds
}

func (pg Pg) GetAll{{$mds}}(tx *Tx) ({{$pkg}}.{{$mds}}, error) {
	mds, err := GetAll{{$mds}}(tx)
	return mds.ServiceData(), err
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

func (pg Pg) Get{{$md}}ById(tx *Tx, id uint64) ({{$pkg}}.{{$md}}, error) {
	md, err := Get{{$md}}ById(tx, id)
	return md.ServiceData(), err
}

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

func (pg Pg) Get{{$mds}}ByComs(tx *Tx, coms []uint64) ({{$pkg}}.{{$mds}}, error) {
	mds, err := Get{{$mds}}ByComs(tx, coms)
	return mds.ServiceData(), err
}

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

func (pg Pg) Count{{$md}}ByUnique(tx *Tx, u {{$pkg}}.{{$md}}Unique) (uint64, error) {
	var by *{{$md}}Unique
	by.SetData(u)
	i, err := Count{{$md}}ByUnique(tx, *by)
	return uint64(i), err
}

func Count{{$md}}ByUnique(tx *Tx, c {{$md}}Unique) (int64, error) {

	n, err := tx.Where("com = ? AND version = ?", c.Com, c.Version).Count(&{{$md}}{})
	if err != nil {
		logger.Error(err)
		return n, err
	}
	return n, nil
}

func (pg Pg) Insert{{$md}}(tx *Tx, mds ...*{{$pkg}}.{{$md}}) error {
	var bys []*{{$md}}
	for i := range mds {
		var by *{{$md}}
		by.SetData(*mds[i])
		bys = append(bys, by)
	}
	return Insert{{$md}}(tx, bys...)
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

func (pg Pg) Update{{$md}}(tx *Tx, md {{$pkg}}.{{$md}}) error {
	var by *{{$md}}
	by.SetData(md)
	return Update{{$md}}(tx, *by)
}

func Update{{$md}}(tx *Tx, md {{$md}}) error {
	if _, err := tx.Where("id = ?", md.ID).Update(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) Delete{{$md}}(tx *Tx, md {{$pkg}}.{{$md}}) error {
	var by *{{$md}}
	by.SetData(md)
	return Delete{{$md}}(tx, *by)
}

func Delete{{$md}}(tx *Tx, md {{$md}}) error {

	if _, err := tx.Delete(md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) Delete{{$md}}ById(tx *Tx, id uint64) error {
	return Delete{{$md}}ById(tx, id)
}

func Delete{{$md}}ById(tx *Tx, id uint64) error {
	var md {{$md}}
	if _, err := tx.Where("id = ?", id).Delete(&md); err != nil {
		logger.Error(err)
		return err
	}
	return nil
}

func (pg Pg) Delete{{$md}}ByUnique(tx *Tx, md {{$pkg}}.{{$md}}Unique) error {
	var by *{{$md}}Unique
	by.SetData(md)
	return Delete{{$md}}ByUnique(tx, *by)
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
