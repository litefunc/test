package pgsql

import (
	"reflect"
	"regexp"
	"strings"
)

func GetTable(md interface{}) string {
	if md == nil {
		return ""
	}

	var t reflect.Type

	val := reflect.ValueOf(md)
	// logger.Debug(val.Kind(), val.Kind() == reflect.Interface, md == nil)

	if val.Kind() == reflect.Ptr {
		t = reflect.TypeOf(md).Elem()
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}
	} else {
		t = reflect.TypeOf(md)
		if t.Kind() == reflect.Slice {
			t = reflect.TypeOf(md).Elem()
		}
	}

	for i := 0; i < t.NumField(); i++ {
		if strings.ToLower(t.Field(i).Name) == "tablename" {
			return t.Field(i).Tag.Get("sql")
		}
	}
	return toSnakeCase(t.Name())
}

func GetColsString(md interface{}) string {
	return strings.Join(GetCols(md), ", ")
}

func GetCols(md interface{}) []string {
	if md == nil {
		return nil
	}

	tb := GetTable(md)

	kindOfJ := reflect.ValueOf(md).Kind()
	if kindOfJ == reflect.Ptr {
		t := reflect.TypeOf(md).Elem()
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		var cols []string
		return getCols(t, tb, cols)
	}

	t := reflect.TypeOf(md)
	if t.Kind() == reflect.Slice {
		t = reflect.TypeOf(md).Elem()
	}
	var cols []string

	return getCols(t, tb, cols)
}

func getCols(t reflect.Type, tb string, cols []string) []string {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column := field.Tag.Get("sql")

		if column == tb {
			continue
		}

		if column == "embed" {
			cols = getCols(field.Type, tb, cols)
			continue
		}

		if column != "" {

			strs := strings.Split(column, ",")
			if strs[0] != "" {
				cols = append(cols, strs[0])
			} else {
				cols = append(cols, toSnakeCase(t.Field(i).Name))
			}

		} else {
			cols = append(cols, toSnakeCase(t.Field(i).Name))
		}
	}
	return cols
}

func toSnakeCase(str string) string {

	matchFirstCap := regexp.MustCompile("(.)([A-Z][a-z]+)")
	matchAllCap := regexp.MustCompile("([a-z0-9])([A-Z])")

	snake := matchFirstCap.ReplaceAllString(str, "${1}_${2}")
	snake = matchAllCap.ReplaceAllString(snake, "${1}_${2}")
	return strings.ToLower(snake)
}
