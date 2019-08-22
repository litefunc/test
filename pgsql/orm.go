package pgsql

import (
	"log"
	"reflect"
	"regexp"
	"strings"
)

func GetTable(md interface{}) string {
	if md == nil {
		return ""
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

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
			return t.Field(i).Tag.Get("db")
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

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	tb := GetTable(md)
	t := reflect.TypeOf(md)
	v := reflect.ValueOf(md)

	kindOfJ := v.Kind()
	if kindOfJ == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Slice {
			t = t.Elem()
		}
		var cols []string
		return getCols(t, tb, cols)
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
	}
	var cols []string

	return getCols(t, tb, cols)
}

func getCols(t reflect.Type, tb string, cols []string) []string {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column := field.Tag.Get("db")

		if column == tb {
			continue
		}

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
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

func GetValues(md interface{}) []interface{} {
	if md == nil {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	tb := GetTable(md)

	t := reflect.TypeOf(md)
	v := reflect.ValueOf(md)
	kindOfJ := v.Kind()
	if kindOfJ == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Slice {
			t = t.Elem()
			v = v.Elem()
		}
		var cols []interface{}
		return getValues(t, v, tb, cols)
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		v = v.Elem()
	}
	var cols []interface{}

	return getValues(t, v, tb, cols)
}

func getValues(t reflect.Type, v reflect.Value, tb string, cols []interface{}) []interface{} {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column := field.Tag.Get("db")

		if column == tb {
			continue
		}

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			cols = getValues(field.Type, v.Field(i), tb, cols)
			continue
		}

		cols = append(cols, v.Field(i).Interface())

	}
	return cols
}

func GetColsValues(md interface{}) map[string]interface{} {
	if md == nil {
		return nil
	}

	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
		}
	}()

	tb := GetTable(md)

	t := reflect.TypeOf(md)
	v := reflect.ValueOf(md)
	kindOfJ := v.Kind()
	if kindOfJ == reflect.Ptr {
		t = t.Elem()
		if t.Kind() == reflect.Slice {
			t = t.Elem()
			v = v.Elem()
		}
		cols := make(map[string]interface{})
		return getColsValues(t, v, tb, cols)
	}

	if t.Kind() == reflect.Slice {
		t = t.Elem()
		v = v.Elem()
	}
	cols := make(map[string]interface{})

	return getColsValues(t, v, tb, cols)
}

func getColsValues(t reflect.Type, v reflect.Value, tb string, cols map[string]interface{}) map[string]interface{} {
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		column := field.Tag.Get("db")

		if column == tb {
			continue
		}

		if field.Type.Kind() == reflect.Struct && field.Anonymous {
			cols = getColsValues(field.Type, v.Field(i), tb, cols)
			continue
		}
		cols[toSnakeCase(t.Field(i).Name)] = v.Field(i).Interface()

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
