package aes

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/hex"
	"io"
	"test/logger"
)

// =====================================================
// string --> bytes --> encrypt --> bytes --> encode --> string
// string <-- bytes <-- decrypt <-- bytes <-- decode <-- string
// =====================================================

func NewKey() (string, error) {
	bytes := make([]byte, 16) //generate a random 16 byte key for AES-128
	if _, err := rand.Read(bytes); err != nil {
		logger.Error(err)
		return "", err
	}

	return hex.EncodeToString(bytes), nil
}

func decodeString(s string) ([]byte, error) {

	b, err := hex.DecodeString(s)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	return b, nil
}

func Encrypt(data, key string) (string, error) {

	by, err := decodeString(key)
	if err != nil {
		return "", err
	}
	b, err := encrypt([]byte(data), by)
	if err != nil {
		return "", err
	}

	return hex.EncodeToString(b), nil
}

func encrypt(data, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	ciphertext := gcm.Seal(nonce, nonce, data, nil)
	return ciphertext, nil
}

func Decrypt(data, key string) (string, error) {

	by0, err := decodeString(data)
	if err != nil {
		return "", err
	}

	by1, err := decodeString(key)
	if err != nil {
		return "", err
	}
	b, err := decrypt(by0, by1)
	if err != nil {
		return "", err
	}
	return string(b), nil
}

func decrypt(data, key []byte) ([]byte, error) {

	block, err := aes.NewCipher(key)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	nonceSize := gcm.NonceSize()
	nonce, ciphertext := data[:nonceSize], data[nonceSize:]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	return plaintext, nil
}
