package main

import (
	_ "embed"
	"errors"
	"fmt"
	"net/http"
	"os"

	"github.com/carsonoid/talk-lets-auth-with-go/pkg/simplejwt"
)

func main() {
	if len(os.Args) != 2 {
		fmt.Printf("USAGE %s <private-ed-key-path>\n", os.Args[0])
		os.Exit(1)
	}

	issuer, err := simplejwt.NewIssuer(os.Args[1])
	if err != nil {
		panic(err)
	}

	auth, err := NewAuthService(issuer)
	if err != nil {
		panic(err)
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/login", auth.HandleLogin)

	fmt.Println("Listening on :8081")
	err = http.ListenAndServe(":8081", mux)
	if err != nil {
		panic(err)
	}
}

// AuthService handles authentication and issues tokens
type AuthService struct {
	issuer *simplejwt.Issuer
}

// NewAuthService creates a new service using the given issuer
func NewAuthService(issuer *simplejwt.Issuer) (*AuthService, error) {
	if issuer == nil {
		return nil, errors.New("issuer is required")

	}

	return &AuthService{
		issuer: issuer,
	}, nil
}

func (a *AuthService) HandleLogin(w http.ResponseWriter, r *http.Request) {
	// check basic auth
	user, pass, ok := r.BasicAuth()
	if !ok {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("missing basic auth")) //nolint
		return
	}

	// This is a bad idea in anything real. But this isn't an auth methods talk
	// this talk is about JWTs so we only have trivial auth checking here
	if user != "admin" || pass != "pass" {
		w.WriteHeader(http.StatusUnauthorized)
		w.Write([]byte("invalid credentials")) //nolint
		return
	}

	tokenString, err := a.issuer.IssueToken("admin", []string{"admin", "basic"})
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("unable to issue token:" + err.Error())) //nolint
		return
	}

	_, _ = w.Write([]byte(tokenString + "\n"))
}
