package bcrypto

import (
	"crypto/rand"

	"golang.org/x/crypto/argon2"
)

// Generate random keys for cryptographic use
func GenKey(KeySize int) ([]byte, error) {
	key := make([]byte, KeySize) // create an array of bytes
	_, err := rand.Read(key)     // fill the array randomly generated values
	if err != nil {
		return nil, err
	}

	return key, nil // return the key
}

// Derivate a key or a password with argon2
func DeriveKey(key []byte) []byte {
	// Generating derivated key from a byte array
	// this function use the recomanded settings for
	// generating a 32 bytes key.
	key = argon2.Key(key, nil, 3, 32*1024, 4, 32)
	return key // return key
}
