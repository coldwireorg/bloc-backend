package tokens

import (
	"crypto/ed25519"
	"crypto/rand"
	"log"
	"time"

	"github.com/kataras/jwt"
)

// private signing key
var JWTPrivateKey ed25519.PrivateKey

// Payload of a JWT token
type Token struct {
	Username   string `json:"username"`
	PrivateKey string `json:"privateKey"`
}

// Generate ed25519 signing keys
func GenerateKeys() {
	// Generate private key, public key can be get with PrivateKey.Public()
	_, PrivateKey, err := ed25519.GenerateKey(rand.Reader)
	if err != nil {
		log.Fatal(err)
	}

	// set the global private key variable with one generated
	JWTPrivateKey = PrivateKey
}

// Generate token:
//	In this function we just put the content of the token
//	like the username, the user id and the expiration time
//	of the token
func Generate(username string, pvKey string, exp time.Duration) string {
	// Define the body
	body := Token{
		Username:   username,
		PrivateKey: pvKey,
	}

	// Define the header
	header := jwt.Claims{
		Expiry:   time.Now().Add(exp).Unix(),
		IssuedAt: time.Now().Unix(),
		Issuer:   "coldwire",
	}

	// Sign the token with a ed25519 private key
	t, err := jwt.Sign(jwt.EdDSA, JWTPrivateKey, body, header)
	if err != nil {
		log.Fatal(err)
	}

	// return token as a string
	return string(t)
}

// Verify JWT tokens:
// 	This simply work by verifying the signature with
//	a ed25519 public key
func Verify(token string) (*jwt.VerifiedToken, error) {
	t := []byte(token) // get token

	// Verify token
	verifiedToken, err := jwt.Verify(jwt.EdDSA, JWTPrivateKey.Public(), t)
	if err != nil {
		return nil, err
	}

	// return token
	return verifiedToken, nil
}

// This function is used for decoding JWT token
// WARNING: this is actually not verifying it!!
func Parse(token string) (Token, error) {
	t, err := jwt.Decode([]byte(token)) // get token from the cooki
	if err != nil {
		return Token{}, err // return void token payload with the error
	}

	// Get token payload content
	var payload Token
	err = t.Claims(&payload)
	if err != nil {
		return Token{}, err
	}

	// return the token payload without errors
	return payload, nil
}
