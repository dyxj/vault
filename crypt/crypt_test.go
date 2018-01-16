package crypt_test

import (
	"reflect"
	"testing"
	"vault/crypt"
)

func TestEncryptDecrypt(t *testing.T) {
	key := "1234"
	payload := "Hello"
	keyE := crypt.HashTo32Bytes([]byte(key))
	encArr, err := crypt.EncryptBytes([]byte(payload), keyE)
	if err != nil {
		t.Fatalf("Error encrypting file: %v", err)
	}
	decArr, err := crypt.DecryptBytes([]byte(encArr), keyE)
	if err != nil {
		t.Fatalf("Error decrypting file: %v", err)
	}
	if !reflect.DeepEqual(decArr, []byte(payload)) {
		t.Fatalf("payload and output does not match")
	}
}

func TestInvalidKey(t *testing.T) {
	expectedError := "cipher: message authentication failed"
	key := "1234"
	payload := "Hello"
	keyE := crypt.HashTo32Bytes([]byte(key))
	invalidKey := crypt.HashTo32Bytes([]byte("123"))
	encArr, err := crypt.EncryptBytes([]byte(payload), keyE)
	if err != nil {
		t.Fatalf("Error encrypting file: %v", err)
	}
	_, err = crypt.DecryptBytes([]byte(encArr), invalidKey)
	if err.Error() != expectedError {
		t.Fatalf("Unexpected error during decryption: %v", err)
	}
}
