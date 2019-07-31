package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"fmt"
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
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	mds := MsFws{md1, md2}

	t := reflect.TypeOf(md1)
	logger.Debug(t.Kind() == reflect.Struct)
	logger.Debug(t.String())

	t1 := reflect.TypeOf(mds)
	logger.Debug(t1.Kind() == reflect.Slice)
	logger.Debug(t1.String())

	field1 := t.Field(1)
	tag1 := field1.Tag.Get("sql")
	logger.Debug(tag1)

	for i := 0; i < t.NumField(); i++ {
		logger.Debug(t.Field(i))
		logger.Debug(t.Field(i).Type)
		p := reflect.New(t.Field(i).Type)
		logger.Debug(p)
	}

	InspectStruct(&md1)

}

func InspectStructV(val reflect.Value) {
	if val.Kind() == reflect.Interface && !val.IsNil() {
		elm := val.Elem()
		if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
			val = elm
		}
	}
	if val.Kind() == reflect.Ptr {
		val = val.Elem()
	}

	for i := 0; i < val.NumField(); i++ {
		valueField := val.Field(i)
		typeField := val.Type().Field(i)
		address := "not-addressable"

		if valueField.Kind() == reflect.Interface && !valueField.IsNil() {
			elm := valueField.Elem()
			if elm.Kind() == reflect.Ptr && !elm.IsNil() && elm.Elem().Kind() == reflect.Ptr {
				valueField = elm
			}
		}

		if valueField.Kind() == reflect.Ptr {
			valueField = valueField.Elem()

		}
		if valueField.CanAddr() {
			address = fmt.Sprintf("0x%X", valueField.Addr().Pointer())
		}

		fmt.Printf("Field Name: %s,\t Field Value: %v,\t Address: %v\t, Field type: %v\t, Field kind: %v\n", typeField.Name,
			valueField.Interface(), address, typeField.Type, valueField.Kind())

		if valueField.Kind() == reflect.Struct {
			InspectStructV(valueField)
		}
	}
}

func InspectStruct(v interface{}) {
	InspectStructV(reflect.ValueOf(v))
}
