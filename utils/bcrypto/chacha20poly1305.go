package bcrypto

import (
	"crypto/rand"
	"errors"

	"golang.org/x/crypto/chacha20poly1305"
)

// XChacha20-Poly1305 encryption function
func Encrypt(crypt []byte, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	// Select a random nonce, and leave capacity for the ciphertext.
	nonce := make([]byte, aead.NonceSize(), aead.NonceSize()+len(crypt)+aead.Overhead())
	if _, err := rand.Read(nonce); err != nil {
		return nil, err
	}

	// Encrypt the message and append the ciphertext to the nonce.
	return aead.Seal(nonce, nonce, crypt, nil), nil
}

// XChacha20-Poly1305 decryption function
func Decrypt(crypted []byte, key []byte) ([]byte, error) {
	aead, err := chacha20poly1305.NewX(key)
	if err != nil {
		return nil, err
	}

	if len(crypted) < aead.NonceSize() {
		return nil, errors.New("ciphertext too short")
	}

	// Split nonce and ciphertext.
	nonce, ciphertext := crypted[:aead.NonceSize()], crypted[aead.NonceSize():]

	// Decrypt the message and check it wasn't tampered with.
	decrypt, err := aead.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return decrypt, nil
}
