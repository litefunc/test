package internal

import (
	"errors"
	"fmt"
	"mstore/logger"
	"time"

	"github.com/dgrijalva/jwt-go"
)

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

func Duplicate() {
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

func Gen(id, secret string, exp time.Duration) string {
	now := time.Now().UTC()
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.StandardClaims{
		Id:        id,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.Add(exp).Unix(),
	})

	mySigningKey := []byte(secret)
	t, err := token.SignedString(mySigningKey)
	if err != nil {
		logger.Panic(err)
	}
	return t
}

func Validate(id uint64, secret, tk string) error {
	// keyFunc return JWT secret key, whitch is used to validate token string
	keyFunc := func(token *jwt.Token) (interface{}, error) {

		c, ok := token.Claims.(*jwt.StandardClaims)
		if !ok {
			logger.Error(ok)
			return nil, errors.New("invalid JWT claim")
		}
		logger.Debug(c)

		return []byte(secret), nil
	}

	token, err := jwt.ParseWithClaims(tk, &jwt.StandardClaims{}, keyFunc)
	if err != nil || !token.Valid {
		logger.Error(id, err, token.Valid)
		return err
	}

	return nil
}
