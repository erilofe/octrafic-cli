package auth

import (
	"fmt"
	"net/http"
	"strings"
)

// BasicAuth represents HTTP Basic authentication
type BasicAuth struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewBasicAuth creates a new Basic authentication provider
func NewBasicAuth(username, password string) *BasicAuth {
	return &BasicAuth{
		Username: username,
		Password: password,
	}
}

// Apply adds Basic authentication to the request
func (b *BasicAuth) Apply(req *http.Request) error {
	if err := b.Validate(); err != nil {
		return err
	}
	req.SetBasicAuth(b.Username, b.Password)
	return nil
}

// Type returns the authentication type
func (b *BasicAuth) Type() string {
	return "basic"
}

// Validate checks if username and password are present
func (b *BasicAuth) Validate() error {
	if strings.TrimSpace(b.Username) == "" {
		return fmt.Errorf("username cannot be empty")
	}
	if strings.TrimSpace(b.Password) == "" {
		return fmt.Errorf("password cannot be empty")
	}
	return nil
}

// Redact returns a copy with password redacted
func (b *BasicAuth) Redact() AuthProvider {
	return &BasicAuth{
		Username: b.Username,
		Password: "***",
	}
}

// String returns a human-readable representation
func (b *BasicAuth) String() string {
	return fmt.Sprintf("Basic Auth (%s)", b.Username)
}
