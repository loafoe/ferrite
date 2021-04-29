package types

import (
	"crypto/x509"
	"encoding/pem"
	"fmt"
	"time"
)

type Cluster struct {
	ID         string    `json:"id" gorm:"primaryKey"`
	PrivateKey string    `json:"private_key,omitempty"`
	PublicKey  string    `json:"public_key,omitempty" gorm:"-"`
	CreatedAt  time.Time `json:"created_at"`
	UpdatedAt  time.Time `json:"updated_at"`
}

func (c Cluster) PublicKeyPEM() (string, error) {
	rsaBlock, rest := pem.Decode([]byte(c.PrivateKey))
	if rsaBlock == nil {
		return "", fmt.Errorf("error decoding: len(rest)=%d", len(rest))
	}
	rsaKey, err := x509.ParsePKCS1PrivateKey(rsaBlock.Bytes)
	if err != nil {
		return "", err
	}
	publicKeyDer := x509.MarshalPKCS1PublicKey(&rsaKey.PublicKey)
	pubKeyBlock := pem.Block{
		Type:    "PUBLIC KEY",
		Headers: nil,
		Bytes:   publicKeyDer,
	}
	pubKeyPem := string(pem.EncodeToMemory(&pubKeyBlock))
	return pubKeyPem, nil
}
