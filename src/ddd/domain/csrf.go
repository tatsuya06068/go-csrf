package csrf

import (
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"time"
)

type Token struct {
	Value     string
	ExpiresAt time.Time
}

// CSRFTokenService is the domain service for generating and validating CSRF tokens.
type CSRFTokenService struct {
	secretKey string
}

func NewCSRFTokenService(secretKey string) *CSRFTokenService {
	return &CSRFTokenService{
		secretKey: secretKey,
	}
}

// GenerateToken generates a new CSRF token.
func (s *CSRFTokenService) GenerateToken() (Token, error) {
	tokenBytes := make([]byte, 32)
	_, err := rand.Read(tokenBytes)
	if err != nil {
		return Token{}, err
	}

	token := base64.StdEncoding.EncodeToString(tokenBytes)
	expiresAt := time.Now().Add(24 * time.Hour) // トークンの有効期限

	return Token{
		Value:     token,
		ExpiresAt: expiresAt,
	}, nil
}

// ValidateToken checks if the provided token is valid.
func (s *CSRFTokenService) ValidateToken(providedToken string, storedToken Token) error {
	if providedToken != storedToken.Value {
		return errors.New("invalid CSRF token")
	}
	if time.Now().After(storedToken.ExpiresAt) {
		return errors.New("CSRF token has expired")
	}
	return nil
}

type CSRFService interface {
	GenerateToken() (*CSRFToken, error)
	ValidateToken(token string) error
}