package simplejwt

import (
	"context"
	"errors"

	"github.com/golang-jwt/jwt"
)

// middlewareContextKey is our custom type to ensure our context values are unique
type middlewareContextKey string

// tokenContextKey is the key used for a parsed token
const tokenContextKey middlewareContextKey = "token"

// ContextWithToken adds the given token to the given context
func ContextWithToken(ctx context.Context, token *jwt.Token) context.Context {
	return context.WithValue(ctx, tokenContextKey, token)
}

// ContextGetToken tries to get the token from the context
// it returns an error if the token is missing or invalid
//
// It DOES NOT validate the token claims or signature
// That would require the public key and should have been handled
// by the process that set the token originally
func ContextGetToken(ctx context.Context) (*jwt.Token, error) {
	val := ctx.Value(tokenContextKey)
	if val == nil {
		return nil, errors.New("no token in context")
	}

	t, ok := val.(*jwt.Token)
	if !ok {
		return nil, errors.New("unexpected token type in context")
	}

	return t, nil
}

// MustContextGetToken parses the token out of the context
// it will panic if the token is not found
func MustContextGetToken(ctx context.Context) *jwt.Token {
	t, err := ContextGetToken(ctx)
	if err != nil {
		panic(err)
	}

	return t
}
