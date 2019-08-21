package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"encoding/json"
	"fmt"
	"time"
)

type T1 struct {
	Time string `json:"time`
}

type T2 struct {
	Time time.Time `json:"time`
}

type T3 struct {
	Time null.Time `json:"time`
}

func main() {

	layout := "2006-01-02"
	str := "2014-11-12"
	t, err := time.Parse(layout, str)

	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(t.Format(layout))

	now := time.Now().UTC()
	logger.Debug(now.Format(time.RFC3339))
	logger.Debug(now.Unix())
	logger.Debug(now.UnixNano())

	var n null.Time
	logger.Debug(n)
	logger.Debug(!n.Valid)
	logger.Debug(n.Time.Unix() == -62135596800)
	logger.Debug(uint64(n.Time.Unix()) == 18446744011573954816)

	var tt time.Time
	logger.Debug(tt)
	logger.Debug(tt.Unix() == -62135596800)
	logger.Debug(uint64(tt.Unix()) == 18446744011573954816)

	by, _ := json.Marshal(T2{Time: now})
	logger.Debug(string(by))
	logger.Debug(now.String())


	t1 := T1{Time: "2019-08-19T01:25:24.781Z"}
	by1, _ := json.Marshal(t1)
	logger.Debug(string(by1))

	var t2 T2
	json.Unmarshal(by, &t2)
	by2, _ := json.Marshal(t2)
	logger.Debug(string(by2))

	var t3 T3
	json.Unmarshal(by, &t3)
	by3, _ := json.Marshal(t3)
	logger.Debug(string(by3))

}
