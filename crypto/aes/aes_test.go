package aes

import (
	"fmt"
	"testing"
)

func TestAes(t *testing.T) {
	const key = "66b41f7a5c3d7a73c20dd79ef94f3ddf"

	if err := eq("abcdef", key); err != nil {
		t.Fatal(err)
	}

	k, err := NewKey()
	if err != nil {
		t.Fatal(err)
	}

	if err := eq("ghij", k); err != nil {
		t.Fatal(err)
	}

}

func eq(data, key string) error {
	s1, err := Encrypt(data, key)
	if err != nil {
		return err
	}

	s2, err := Decrypt(s1, key)
	if err != nil {
		return err
	}
	if s2 != data {
		return fmt.Errorf(`want:%s, got:%s`, data, s2)
	}
	return nil
}
