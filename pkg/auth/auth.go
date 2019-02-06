package auth

import (
	"context"
	"errors"
	"time"

	jwt "github.com/dgrijalva/jwt-go"
	kitjwt "github.com/go-kit/kit/auth/jwt"
)

var (
	// ErrUnknowCredentials raised when credentials are not found
	ErrUnknowCredentials = errors.New("unknown credentials")
	// ErrInvalidToken raised when token is not valid
	ErrInvalidToken = errors.New("invalid token")
	// ErrInvalidClaims ...
	ErrInvalidClaims = errors.New("invalid claims")
	// ErrClaimsSubMissing ...
	ErrClaimsSubMissing = errors.New("cannot find claim 'sub' in token")
	// ErrJWTNotFound ...
	ErrJWTNotFound = errors.New("cannot retrieve token from context")
)

type auth struct {
	tokenDuration time.Duration
	signingKey    []byte
	keyFunc       jwt.Keyfunc
}

// Service ...
type Service interface {
	// Generate a token
	GetJWT(id string) (string, error)

	// Get the sub from the token
	ExtractSubFromContext(ctx context.Context) (string, error)
}

// NewService ...
func NewService(tokendDuration time.Duration, key []byte, keyFunc jwt.Keyfunc) Service {
	return auth{
		tokenDuration: tokendDuration,
		signingKey:    key,
		keyFunc:       keyFunc,
	}
}

// GetJWT returns a token
func (a auth) GetJWT(id string) (string, error) {
	claims := jwt.StandardClaims{
		Subject:   id,
		ExpiresAt: time.Now().Add(time.Second * a.tokenDuration).Unix(),
		IssuedAt:  jwt.TimeFunc().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString(a.signingKey)
}

// ExtractUserID retrieves the user id from the token
func (a auth) ExtractSubFromContext(ctx context.Context) (string, error) {
	ownertoken := ctx.Value(kitjwt.JWTTokenContextKey)
	tokenString, ok := ownertoken.(string)
	if !ok {
		return "", ErrJWTNotFound
	}

	token, err := jwt.Parse(tokenString, a.keyFunc)
	if err != nil {
		return "", err
	}

	if !token.Valid {
		return "", ErrInvalidToken
	}

	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		return "", ErrInvalidClaims
	}

	sub := ""
	if s, ok := claims["sub"]; ok {
		sub = s.(string)
	} else {
		return "", ErrClaimsSubMissing
	}

	return sub, nil
}
