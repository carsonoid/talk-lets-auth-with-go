package simplejwt

import (
	"crypto"
	"fmt"
	"io/ioutil"

	"github.com/golang-jwt/jwt"
)

// START VALIDATOR OMIT
// Validator does parsing and validation of JWT token
type Validator struct {
	key crypto.PublicKey
}

// NewValidator returns a new validator by parsing the given file path as a ed25519 public key
func NewValidator(publicKeyPath string) (*Validator, error) {
	keyBytes, err := ioutil.ReadFile(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to read public key file: %w", err)
	}

	key, err := jwt.ParseEdPublicKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	return &Validator{
		key: key,
	}, nil
}

// END VALIDATOR OMIT

// GetToken attempts to get a token from the given string
// it validates both the signature and claim and returns nil and an err if invalid
func (v *Validator) GetToken(tokenString string) (*jwt.Token, error) {
	// jwt.Parse also does signature verify and claim validation
	token, err := jwt.Parse(
		tokenString,
		// the func below is to help figure out if the token came from a key we trust
		// our implementation assumes a single trusted private key
		//
		// NOTE: this is where you would handle key rotation or multiple trusted issuers
		func(token *jwt.Token) (interface{}, error) {
			// Check to see if the token uses the expected signing method
			if _, ok := token.Method.(*jwt.SigningMethodEd25519); !ok {
				return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
			}

			// return the single public key we trust
			return v.key, nil
		})
	if err != nil {
		return nil, fmt.Errorf("unable to parse token string: %w", err)
	}

	return token, nil
}
