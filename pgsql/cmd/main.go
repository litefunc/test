package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"test/pgsql"
	"time"

	"github.com/lib/pq"
)

type MsFws []MsFw

type MsFw struct {
	TableName  struct{} `json:"-" sql:"cloud.msfw"`
	ID         uint64   `json:"id" sql:",pk"`
	MsFwUnique `sql:"extent"`
	Bucket     string      `json:"bucket"`
	Obj        string      `json:"obj"`
	Time       time.Time   `json:"time"`
	Tag        null.String `json:"tag"`
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

var matchFirstCap = regexp.MustCompile("(.)([A-Z][a-z]+)")
var matchAllCap = regexp.MustCompile("([a-z0-9])([A-Z])")

func ToSnakeCase(str string) string {
	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}

func getTable(md interface{}) string {
	t := reflect.TypeOf(md)
	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Field(i).Name) == "tablename" {
			return t.Field(i).Tag.Get("sql")
		}
	}
	return ToSnakeCase(t.Name())
}

func getCols(md interface{}) []string {

	tb := getTable(md)
	t := reflect.TypeOf(md)
	var cols []string
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		c := strings.Replace(field.Tag.Get("sql"), " ", "", -1)
		if c == tb {
			continue
		}

		if c != "" {

			strs := strings.Split(c, ",")
			if strs[0] != "" {
				cols = append(cols, strs[0])
			} else {
				cols = append(cols, ToSnakeCase(t.Field(i).Name))
			}

		} else {
			cols = append(cols, ToSnakeCase(t.Field(i).Name))
		}

	}
	return cols
}

func printTags(t reflect.Type, tb string) {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column := field.Tag.Get("sql")

		if column == tb {
			continue
		}

		if column == "extent" {
			printTags(field.Type, tb)
			continue
		}

		if column != "" {

			strs := strings.Split(column, ",")
			if strs[0] != "" {
				fmt.Println("column: ", strs[0])
			} else {
				fmt.Println("column: ", ToSnakeCase(t.Field(i).Name))
			}

		} else {
			fmt.Println("column: ", ToSnakeCase(t.Field(i).Name))
		}
	}
}

// func printTags(t reflect.Type) {
// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)

// 		if field.Type.Kind() == reflect.Struct {
// 			printTags(field.Type)
// 			continue
// 		}

// 		column := field.Tag.Get("sql")

// 		if column != "" {

// 			strs := strings.Split(column, ",")
// 			if strs[0] != "" {
// 				fmt.Println("column: ", strs[0])
// 			} else {
// 				fmt.Println("column: ", ToSnakeCase(t.Field(i).Name))
// 			}

// 		} else {
// 			fmt.Println("column: ", ToSnakeCase(t.Field(i).Name))
// 		}
// 	}
// }

func main() {

	var db pgsql.DB
	db.Select("id", "com", "ver").From("cloud.msfw").Where("com=? AND gp IN(?)", 1, pq.Array([]uint64{1, 2, 3})).Order("com DESC", "gp ASC").Limit(10).SQL()

	md1 := NewMsFw("v1", 1, "bk1", "obj1", time.Now().UTC(), "tag1")
	md2 := NewMsFw("v2", 2, "bk2", "obj2", time.Now().UTC(), "tag2")
	mds := MsFws{md1, md2}

	t := reflect.TypeOf(md1)
	logger.Debug(t.Kind() == reflect.Struct)
	logger.Debug(t.String())
	logger.Debug(t.Name())

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

	// var u MsFwUnique
	// logger.Debug(getTable(md1))
	// logger.Debug(getTable(u))
	logger.Debug(getCols(md1))
	// logger.Debug(getCols(u))
	// printTags(reflect.TypeOf(&md1).Elem())
	printTags(reflect.TypeOf(md1), "cloud.msfw")
}
