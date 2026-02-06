package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// BearerAuth represents Bearer token authentication
type BearerAuth struct {
	Token string `json:"token"`
}

// NewBearerAuth creates a new Bearer authentication provider
func NewBearerAuth(token string) *BearerAuth {
	return &BearerAuth{Token: token}
}

// Apply adds the Bearer token to the Authorization header
func (b *BearerAuth) Apply(req *http.Request) error {
	if err := b.Validate(); err != nil {
		return err
	}
	req.Header.Set("Authorization", "Bearer "+b.Token)
	return nil
}

// Type returns the authentication type
func (b *BearerAuth) Type() string {
	return "bearer"
}

// Validate checks if the token is present
func (b *BearerAuth) Validate() error {
	if strings.TrimSpace(b.Token) == "" {
		return fmt.Errorf("bearer token cannot be empty")
	}
	return nil
}

// Redact returns a copy with the token redacted
func (b *BearerAuth) Redact() AuthProvider {
	return &BearerAuth{
		Token: RedactString(b.Token),
	}
}

// String returns a human-readable representation
func (b *BearerAuth) String() string {
	return fmt.Sprintf("Bearer Token (%s)", RedactString(b.Token))
}
