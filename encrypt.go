package main

import (
	"crypto/rand"
	"crypto/sha256"
	"errors"
	"golang.org/x/crypto/chacha20poly1305"
)

func encrypt(plaintext []byte, pass string) ([]byte, error) {
	key := sha256.Sum256([]byte(pass))
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}

	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(plaintext)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	return aead.Seal(nonce, nonce, plaintext, nil), nil
}

func decrypt(encryptedMsg []byte, pass string) ([]byte, error) {
	key := sha256.Sum256([]byte(pass))
	aead, err := chacha20poly1305.NewX(key[:])
	if err != nil {
		return nil, err
	}

	if len(encryptedMsg) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := encryptedMsg[:aead.NonceSize()], encryptedMsg[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	plaintext, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}
