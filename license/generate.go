package license

import (
	"crypto/ed25519"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"time"
)

// KeyPair holds the public and private keys for license signing.
type KeyPair struct {
	PublicKey  string `json:"public_key"`
	PrivateKey string `json:"private_key"`
}

// GenerateKeyPair generates a new Ed25519 keypair for license signing.
func GenerateKeyPair() (*KeyPair, error) {
	publicKey, privateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		return nil, fmt.Errorf("failed to generate keypair: %w", err)
	}

	return &KeyPair{
		PublicKey:  base64.RawURLEncoding.EncodeToString(publicKey),
		PrivateKey: base64.RawURLEncoding.EncodeToString(privateKey),
	}, nil
}

// GenerateLicenseKey generates a signed license key.
func GenerateLicenseKey(privateKeyBase64, email string, tier Tier, validDays int) (string, error) {
	// Decode private key
	privateKey, err := base64.RawURLEncoding.DecodeString(privateKeyBase64)
	if err != nil {
		return "", fmt.Errorf("invalid private key: %w", err)
	}

	if len(privateKey) != ed25519.PrivateKeySize {
		return "", fmt.Errorf("invalid private key size")
	}

	// Create payload
	now := time.Now()
	payload := Payload{
		Email:     email,
		Tier:      tier,
		IssuedAt:  now.Unix(),
		ExpiresAt: now.AddDate(0, 0, validDays).Unix(),
	}

	payloadBytes, err := json.Marshal(payload)
	if err != nil {
		return "", fmt.Errorf("failed to marshal payload: %w", err)
	}

	payloadBase64 := base64.RawURLEncoding.EncodeToString(payloadBytes)

	// Sign the payload
	signature := ed25519.Sign(privateKey, []byte(payloadBase64))
	signatureBase64 := base64.RawURLEncoding.EncodeToString(signature)

	// Create signed license
	signed := SignedLicense{
		Payload:   payloadBase64,
		Signature: signatureBase64,
	}

	signedBytes, err := json.Marshal(signed)
	if err != nil {
		return "", fmt.Errorf("failed to marshal signed license: %w", err)
	}

	return base64.RawURLEncoding.EncodeToString(signedBytes), nil
}
