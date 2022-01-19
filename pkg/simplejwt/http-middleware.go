package simplejwt

import (
	"fmt"
	"net/http"
	"strings"
)

// Middleware handles all jwt parsing and validation automatically when used
type Middleware struct {
	// embed the validator to make token calls cleaner
	Validator
}

// NewMiddleware creates a new middleware that validates using the
// given public key file
func NewMiddleware(publicKeyPath string) (*Middleware, error) {
	validator, err := NewValidator(publicKeyPath)
	if err != nil {
		return nil, fmt.Errorf("unable to create validator: %w", err)
	}

	return &Middleware{
		Validator: *validator,
	}, nil
}

func (m *Middleware) HandleHTTP(h http.Handler) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		parts := strings.Split(r.Header.Get("Authorization"), " ")
		if len(parts) < 2 || parts[0] != "Bearer" {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("missing or invalid authorization header")) //nolint
			return
		}
		tokenString := parts[1]

		token, err := m.GetToken(tokenString)
		if err != nil {
			w.WriteHeader(http.StatusUnauthorized)
			w.Write([]byte("invalid token: " + err.Error())) //nolint
			return
		}

		// Get a new context with the parsed token
		ctx := ContextWithToken(r.Context(), token)

		fmt.Println("* HTTP SERVER middleware validated and set set token")

		// call the next handler with the updated context
		h.ServeHTTP(w, r.WithContext(ctx))
	}
}
