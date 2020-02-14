package main

import (
	"VodoPlay/logger"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	for i := 0; i < 10; i++ {
		gen()
	}
}

func gen() {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sn":         "N0001",
		"mac":        "ab:01:23:cd:45:67",
		"expiredate": time.Now().Add(time.Hour * time.Duration(1)).Unix(),
	})

	now := time.Now()
	logger.Debug(now, now.Nanosecond())
	mySigningKey := []byte(fmt.Sprintf(`%s%d`, "test", now.Nanosecond()))
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Error(err)
		return
	}
	logger.Debug(t)
}
