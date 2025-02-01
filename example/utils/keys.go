package utils

import (
	"crypto/ed25519"
	"crypto/x509"
	"encoding/pem"
	"errors"
	"os"
)

func LoadPublicKeyFromDisk(path string) (ed25519.PublicKey, error) {
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	block, _ := pem.Decode(raw)
	if block == nil {
		return nil, errors.New("failed to decode PEM block")
	}

	parsed, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return nil, err
	}

	if key, ok := parsed.(ed25519.PublicKey); ok {
		return key, nil
	} else {
		return nil, errors.New("key is not an Ed25519 public key")
	}
}
