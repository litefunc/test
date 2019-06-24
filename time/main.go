package main

import (
	"cloud/lib/logger"
	"cloud/lib/null"
	"fmt"
	"time"
)

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

}
