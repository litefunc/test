package main

import (
	"test/logger"
)

type Eq interface {
	eq(Eq) bool
	val() interface{}
}

type Int int

func (rec Int) val() interface{} {

	return rec
}

func (rec Int) eq(e Eq) bool {
	i, ok := e.val().(Int)
	if !ok {
		return false
	}
	return i == rec
}

// func (rec Int) eq(e Eq) bool {
// 	i, ok := e.(int)
// 	return reflect.DeepEqual(rec, e)
// }

func unique(es ...Eq) []Eq {
	if len(es) == 0 {
		return []Eq{}
	}
	list := []Eq{}
	m := make(map[Eq]bool)
	for _, v := range es {
		if _, ok := m[v]; !ok {
			m[v] = true
			list = append(list, v)
		}
	}

	return list
}

func main() {
	i := Int(1)
	j := Int(1)
	k := Int(2)
	logger.Debug(i.eq(j))
	logger.Debug(i.eq(k))

	logger.Debug(unique(i, j, k))
}
