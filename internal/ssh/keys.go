package ssh

import (
	"crypto/ed25519"
	"crypto/rand"
	"crypto/rsa"
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"os"

	"golang.org/x/crypto/ssh"
)

// GenerateRSAKey generates an RSA private key of the given bit size.
func GenerateRSAKey(bits int) (*rsa.PrivateKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, fmt.Errorf("failed to generate RSA private key: %w", err)
	}
	return privateKey, nil
}

// GenerateEd25519Key generates an Ed25519 private key.
func GenerateEd25519Key() (ed25519.PrivateKey, error) {
	_, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate Ed25519 private key: %w", err)
	}
	return privateKey, nil
}

// WritePrivateKey writes a private key to a file in PEM format.
func WritePrivateKey(key interface{}, path string) error {
	var pemBlock *pem.Block

	switch k := key.(type) {
	case *rsa.PrivateKey:
		pemBlock = &pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: x509.MarshalPKCS1PrivateKey(k),
		}
	case ed25519.PrivateKey:
		b, err := x509.MarshalPKCS8PrivateKey(k)
		if err != nil {
			return fmt.Errorf("unable to marshal Ed25519 private key: %w", err)
		}
		pemBlock = &pem.Block{
			Type:  "PRIVATE KEY", // Generic private key type
			Bytes: b,
		}
	default:
		return fmt.Errorf("unsupported private key type: %T", key)
	}

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0600)
	if err != nil {
		return fmt.Errorf("failed to open private key file: %w", err)
	}
	defer file.Close()

	if err := pem.Encode(file, pemBlock); err != nil {
		return fmt.Errorf("failed to encode private key: %w", err)
	}

	return nil
}

// WritePublicKey writes a public key to a file in OpenSSH authorized_keys format.
func WritePublicKey(key interface{}, path string) error {
	var publicKey ssh.PublicKey
	var err error

	switch k := key.(type) {
	case *rsa.PrivateKey:
		publicKey, err = ssh.NewPublicKey(&k.PublicKey)
	case ed25519.PrivateKey:
		publicKey, err = ssh.NewPublicKey(k.Public())
	default:
		return fmt.Errorf("unsupported public key type: %T", key)
	}

	if err != nil {
		return fmt.Errorf("failed to create public key: %w", err)
	}

	pubKeyBytes := ssh.MarshalAuthorizedKey(publicKey)

	file, err := os.OpenFile(path, os.O_WRONLY|os.O_CREATE|os.O_TRUNC, 0644)
	if err != nil {
		return fmt.Errorf("failed to open public key file: %w", err)
	}
	defer file.Close()

	if _, err := file.Write(pubKeyBytes); err != nil {
		return fmt.Errorf("failed to write public key: %w", err)
	}

	return nil
}