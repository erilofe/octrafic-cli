package auth

import (
	"fmt"
	"net/http"
)

// AuthProvider applies authentication to HTTP requests
type AuthProvider interface {
	// Apply adds authentication to the request
	Apply(req *http.Request) error

	// Type returns the authentication type identifier
	Type() string

	// Validate checks if the configuration is valid
	Validate() error

	// Redact returns a copy with sensitive data hidden (for logging)
	Redact() AuthProvider
}

// NoAuth represents no authentication
type NoAuth struct{}

func (n *NoAuth) Apply(req *http.Request) error {
	return nil
}

func (n *NoAuth) Type() string {
	return "none"
}

func (n *NoAuth) Validate() error {
	return nil
}

func (n *NoAuth) Redact() AuthProvider {
	return n
}

// ParseAuthType converts a string to an auth type validator
func ParseAuthType(authType string) (string, error) {
	validTypes := map[string]bool{
		"none":   true,
		"bearer": true,
		"apikey": true,
		"basic":  true,
	}

	if !validTypes[authType] {
		return "", fmt.Errorf("invalid auth type: %s (valid: none, bearer, apikey, basic)", authType)
	}

	return authType, nil
}

// RedactString hides sensitive data for logging
func RedactString(s string) string {
	if len(s) == 0 {
		return "<empty>"
	}
	if len(s) <= 8 {
		return "***"
	}
	return s[:4] + "***" + s[len(s)-4:]
}
