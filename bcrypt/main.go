package main

import (
	"cloud/lib/logger"

	"golang.org/x/crypto/bcrypt"
)

const (
	pwd  = "abcdef"
	hash = "$2a$10$KNrqToE3Yb/quxJqEnZeb.kLm3YWi3iIR4C..py4Gpl6tVg0VFHJq"
)

func main() {
	by, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
	}
	h := string(by)
	logger.Debug(h)

	if err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(h)); err != nil {
		logger.Error(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(pwd)); err != nil {
		logger.Error(err)
	}

	if err := bcrypt.CompareHashAndPassword([]byte(h), []byte(hash)); err != nil {
		logger.Error(err)
	}

	by1, err := bcrypt.GenerateFromPassword([]byte(pwd), bcrypt.DefaultCost)
	if err != nil {
		logger.Error(err)
	}

	if err := bcrypt.CompareHashAndPassword(by, by1); err != nil {
		logger.Error(err)
	}
}
