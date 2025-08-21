package utils

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"

	"github.com/99designs/keyring"
)

const (
	keyringService = "ssh-manager"
	keyringUser    = "encryption-key"
)

// getEncryptionKey retrieves or generates a 32-byte AES key from the system keyring.
func getEncryptionKey() ([]byte, error) {
	kr, err := keyring.Open(keyring.Config{
		ServiceName: keyringService,
	})
	if err != nil {
		return nil, fmt.Errorf("failed to open keyring: %w", err)
	}

	key, err := kr.Get(keyringUser)
	if err == keyring.ErrKeyNotFound {
		// Key not found, generate a new one
		keyBytes := make([]byte, 32) // AES-256 key
		if _, err := io.ReadFull(rand.Reader, keyBytes); err != nil {
			return nil, fmt.Errorf("failed to generate encryption key: %w", err)
		}

		err = kr.Set(keyring.Item{
			Key:  keyringUser,
			Data: keyBytes,
		})
		if err != nil {
			return nil, fmt.Errorf("failed to save encryption key to keyring: %w", err)
		}
		return keyBytes, nil
	} else if err != nil {
		return nil, fmt.Errorf("failed to get encryption key from keyring: %w", err)
	}

	return key.Data, nil
}

// Encrypt encrypts plaintext using AES-256-GCM.
func Encrypt(plaintext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("failed to generate nonce: %w", err)
	}

	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Decrypt decrypts ciphertext using AES-256-GCM.
func Decrypt(ciphertext string) (string, error) {
	key, err := getEncryptionKey()
	if err != nil {
		return "", err
	}

	data, err := base64.StdEncoding.DecodeString(ciphertext)
	if err != nil {
		return "", fmt.Errorf("failed to decode base64 ciphertext: %w", err)
	}

	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed to create cipher: %w", err)
	}

	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("failed to create GCM: %w", err)
	}

	nonceSize := gcm.NonceSize()
	if len(data) < nonceSize {
		return "", errors.New("ciphertext too short")
	}

	nonce, ciphertextBytes := data[:nonceSize], data[nonceSize:]
	plaintextBytes, err := gcm.Open(nil, nonce, ciphertextBytes, nil)
	if err != nil {
		return "", fmt.Errorf("failed to decrypt: %w", err)
	}

	return string(plaintextBytes), nil
}