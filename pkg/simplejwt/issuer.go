package simplejwt

import (
	"crypto"
	"fmt"
	"io/ioutil"
	"time"

	"github.com/golang-jwt/jwt"
)

// START ISSUER OMIT
// Issuer handles JWT issuing
type Issuer struct {
	key crypto.PrivateKey
}

// NewIssuer creates a new issuer by parsing the given path as a ed25519 private key
func NewIssuer(privateKeyPath string) (*Issuer, error) {
	keyBytes, err := ioutil.ReadFile(privateKeyPath)
	if err != nil {
		panic(fmt.Errorf("unable to read private key file: %w", err))
	}

	key, err := jwt.ParseEdPrivateKeyFromPEM(keyBytes)
	if err != nil {
		return nil, fmt.Errorf("unable to parse as ed private key: %w", err)
	}

	return &Issuer{
		key: key,
	}, nil
}

// END ISSUER OMIT

// IssueToken issues a new token for the given user with the given roles
func (i *Issuer) IssueToken(user string, roles []string) (string, error) {
	now := time.Now()
	token := jwt.NewWithClaims(&jwt.SigningMethodEd25519{}, jwt.MapClaims{
		// standardized claims
		"aud": "api",
		"nbf": now.Unix(),
		"iat": now.Unix(),
		"exp": now.Add(time.Minute).Unix(),
		"iss": "http://localhost:8081",

		// user is custom claim for the validated user
		"user": user,

		// roles is a list of roles attached to the user
		// it shows that claims can have more complex value types
		"roles": roles,
	})

	// Sign and get the complete encoded token as a string using the secret
	tokenString, err := token.SignedString(i.key)
	if err != nil {
		return "", fmt.Errorf("unable to sign token: %w", err)
	}

	return tokenString, nil
}
