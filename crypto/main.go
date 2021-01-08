package main

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"test/logger"
)

func createHash(key string) string {
	hasher := sha256.New()
	hasher.Write([]byte(key))
	logger.Debug(len(hasher.Sum(nil)))
	s := hex.EncodeToString(hasher.Sum(nil))
	logger.Debug(s)
	return s
}

// func createHash(key string) string {
// 	return key
// }

func decodeKey(key string) ([]byte, error) {
	s := createHash(key)
	b, err := hex.DecodeString(s)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	return b, nil
}

func encrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := decodeKey(passphrase)
	if err != nil {
		return []byte{}, err
	}
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

func decrypt(data []byte, passphrase string) ([]byte, error) {
	key, err := decodeKey(passphrase)
	if err != nil {
		return []byte{}, err
	}
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

func encryptFile(filename string, data []byte, passphrase string) error {
	f, err := os.Create(filename)
	if err != nil {
		logger.Error(err)
		return err
	}
	defer f.Close()
	b, err := encrypt(data, passphrase)
	if err != nil {
		return err
	}
	if _, err := f.Write(b); err != nil {
		if err != nil {
			logger.Error(err)
			return err
		}
	}
	return nil
}

func decryptFile(filename string, passphrase string) ([]byte, error) {
	data, err := ioutil.ReadFile(filename)
	if err != nil {
		logger.Error(err)
		return []byte{}, err
	}
	return decrypt(data, passphrase)
}

func main() {

	pwd := "8c6cad0a3e8ec2b6ed255a4eda33ed5e"

	fmt.Println("Starting the application...")
	ciphertext, _ := encrypt([]byte("Hello World"), pwd)
	fmt.Printf("Encrypted: %x\n", ciphertext)
	plaintext, _ := decrypt(ciphertext, pwd)
	fmt.Printf("Decrypted: %s\n", plaintext)
	encryptFile("sample.txt", []byte("Hello World"), "password1")
	b, _ := decryptFile("sample.txt", "password1")
	fmt.Println(string(b))
}
