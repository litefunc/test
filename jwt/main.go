package main

import (
	"VodoPlay/logger"
	"fmt"
	"time"

	"github.com/dgrijalva/jwt-go"
)

func main() {
	now := time.Now().UTC()

	m := make(map[string]bool)
	m1 := make(map[string]bool)
	m2 := make(map[string]bool)

	n := 100

	logger.Info(0)
	for i := 0; i < n; i++ {
		tk := gen(now)
		if duplicate(m, tk) {
			logger.Warn(true)
			break
		}
	}
	logger.Info(1)
	for i := 0; i < n; i++ {
		tk := gen1(now)
		if duplicate(m1, tk) {
			logger.Warn(true)
			break
		}
	}
	logger.Info(2)
	for i := 0; i < n; i++ {

		tk := gen2()
		if duplicate(m2, tk) {
			logger.Warn(true)
			break
		}
	}
}

func duplicate(m map[string]bool, tk string) bool {
	if _, ok := m[tk]; ok {
		return true
	}
	m[tk] = true
	return false
}

func gen(exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sn":         "N0001",
		"mac":        "ab:01:23:cd:45:67",
		"expiredate": exp,
	})

	mySigningKey := []byte(fmt.Sprintf(`%s%d`, "test", time.Now().Nanosecond()))
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Panic(err)
	}
	return t
}

func gen1(exp time.Time) string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sn":         "N0001",
		"mac":        "ab:01:23:cd:45:67",
		"expiredate": exp,
	})

	mySigningKey := []byte(fmt.Sprintf(`%s%d`, "test", 1))
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Panic(err)
	}
	return t
}

func gen2() string {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sn":         "N0001",
		"mac":        "ab:01:23:cd:45:67",
		"expiredate": time.Now().UTC(),
	})

	mySigningKey := []byte(fmt.Sprintf(`%s%d`, "test", 2))
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Panic(err)
	}
	return t
}
