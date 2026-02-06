package auth

import (
	"net/http"
	"testing"
)

func TestBearerAuth(t *testing.T) {
	auth := NewBearerAuth("test-token-123")

	// Test validation
	if err := auth.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	// Test Apply
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	if err := auth.Apply(req); err != nil {
		t.Errorf("Apply failed: %v", err)
	}

	expected := "Bearer test-token-123"
	if got := req.Header.Get("Authorization"); got != expected {
		t.Errorf("expected %q, got %q", expected, got)
	}

	// Test Type
	if auth.Type() != "bearer" {
		t.Errorf("expected type 'bearer', got %q", auth.Type())
	}

	// Test Redact
	redacted := auth.Redact()
	if bearer, ok := redacted.(*BearerAuth); ok {
		if bearer.Token == "test-token-123" {
			t.Error("token was not redacted")
		}
	}
}

func TestAPIKeyAuth(t *testing.T) {
	auth := NewAPIKeyAuth("X-API-Key", "my-secret-key", "header")

	// Test validation
	if err := auth.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	// Test Apply to header
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	if err := auth.Apply(req); err != nil {
		t.Errorf("Apply failed: %v", err)
	}

	if got := req.Header.Get("X-API-Key"); got != "my-secret-key" {
		t.Errorf("expected 'my-secret-key', got %q", got)
	}

	// Test Apply to query
	authQuery := NewAPIKeyAuth("api_key", "my-secret-key", "query")
	req2, _ := http.NewRequest("GET", "http://example.com", nil)
	if err := authQuery.Apply(req2); err != nil {
		t.Errorf("Apply failed: %v", err)
	}

	if got := req2.URL.Query().Get("api_key"); got != "my-secret-key" {
		t.Errorf("expected 'my-secret-key' in query, got %q", got)
	}
}

func TestBasicAuth(t *testing.T) {
	auth := NewBasicAuth("user123", "pass456")

	// Test validation
	if err := auth.Validate(); err != nil {
		t.Errorf("validation failed: %v", err)
	}

	// Test Apply
	req, _ := http.NewRequest("GET", "http://example.com", nil)
	if err := auth.Apply(req); err != nil {
		t.Errorf("Apply failed: %v", err)
	}

	user, pass, ok := req.BasicAuth()
	if !ok {
		t.Error("BasicAuth not set")
	}
	if user != "user123" {
		t.Errorf("expected user 'user123', got %q", user)
	}
	if pass != "pass456" {
		t.Errorf("expected pass 'pass456', got %q", pass)
	}

	// Test Redact
	redacted := auth.Redact()
	if basic, ok := redacted.(*BasicAuth); ok {
		if basic.Password != "***" {
			t.Error("password was not redacted")
		}
		if basic.Username != "user123" {
			t.Error("username should not be redacted")
		}
	}
}

func TestNoAuth(t *testing.T) {
	auth := &NoAuth{}

	req, _ := http.NewRequest("GET", "http://example.com", nil)
	if err := auth.Apply(req); err != nil {
		t.Errorf("NoAuth Apply should not fail: %v", err)
	}

	if auth.Type() != "none" {
		t.Errorf("expected type 'none', got %q", auth.Type())
	}
}

func TestValidation(t *testing.T) {
	// Empty bearer token
	bearer := NewBearerAuth("")
	if err := bearer.Validate(); err == nil {
		t.Error("empty bearer token should fail validation")
	}

	// Empty API key
	apikey := NewAPIKeyAuth("", "value", "header")
	if err := apikey.Validate(); err == nil {
		t.Error("empty API key name should fail validation")
	}

	// Invalid location
	apikey2 := NewAPIKeyAuth("key", "value", "invalid")
	if err := apikey2.Validate(); err == nil {
		t.Error("invalid location should fail validation")
	}

	// Empty username
	basic := NewBasicAuth("", "pass")
	if err := basic.Validate(); err == nil {
		t.Error("empty username should fail validation")
	}
}
