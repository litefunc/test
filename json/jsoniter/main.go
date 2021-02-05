package main

import "test/json/jsoniter/internal"

type ColorGroup struct {
	ID     int
	Name   string
	Colors []string
}

func main() {

	// group := ColorGroup{
	// 	ID:     1,
	// 	Name:   "Reds",
	// 	Colors: []string{"Crimson", "Red", "Ruby", "Maroon"},
	// }
	// // b, err := jsoniter.Marshal(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	// b, err := jsoniter.Marshal(group)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// logger.Debug(string(b))
	// var list []string
	// jsoniter.Get(b, "Colors").ToVal(&list)
	// logger.Debug(list)
	// logger.Debug(jsoniter.Get(b, "Colors", 10).ToString())

	// val := []byte(`{"ID":1,"Name":"Reds","Colors":["Crimson","Red","Ruby","Maroon"]}`)
	// logger.Debug(string(val))
	// logger.Debug(jsoniter.Get(val, "Colors", 0).ToString())
	// m := make(map[string]interface{})
	// if err := jsoniter.Unmarshal(val, &m); err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// logger.Debug(m)
	// s, err := jsoniter.MarshalToString(m)
	// if err != nil {
	// 	logger.Error(err)
	// 	return
	// }
	// logger.Debug(s)

	internal.Null()

}
