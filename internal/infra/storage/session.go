package storage

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"time"
)

const sessionFile = "session.json"

// Session represents a user authentication session
type Session struct {
	AccessToken  string    `json:"access_token"`
	RefreshToken string    `json:"refresh_token"`
	UserEmail    string    `json:"user_email"`
	ExpiresAt    time.Time `json:"expires_at"`
}

// GetSessionPath returns the path to the session file
func GetSessionPath() (string, error) {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return "", fmt.Errorf("failed to get home directory: %w", err)
	}

	storagePath := filepath.Join(homeDir, storageDir)
	if err := os.MkdirAll(storagePath, 0755); err != nil {
		return "", fmt.Errorf("failed to create storage directory: %w", err)
	}

	return filepath.Join(storagePath, sessionFile), nil
}

// LoadSession loads the session from disk
func LoadSession() (*Session, error) {
	path, err := GetSessionPath()
	if err != nil {
		return nil, err
	}

	data, err := os.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("failed to read session file: %w", err)
	}

	var session Session
	if err := json.Unmarshal(data, &session); err != nil {
		return nil, fmt.Errorf("failed to unmarshal session: %w", err)
	}

	return &session, nil
}

// SaveSession saves the session to disk
func SaveSession(session *Session) error {
	path, err := GetSessionPath()
	if err != nil {
		return err
	}

	data, err := json.MarshalIndent(session, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to marshal session: %w", err)
	}

	// Use 0600 for security (only owner can read/write)
	if err := os.WriteFile(path, data, 0600); err != nil {
		return fmt.Errorf("failed to write session file: %w", err)
	}

	return nil
}

// ClearSession removes the session file
func ClearSession() error {
	path, err := GetSessionPath()
	if err != nil {
		return err
	}

	if err := os.Remove(path); err != nil && !os.IsNotExist(err) {
		return fmt.Errorf("failed to remove session file: %w", err)
	}

	return nil
}

// IsExpired checks if the session is expired
func (s *Session) IsExpired() bool {
	return time.Now().After(s.ExpiresAt)
}

// SessionExists checks if a session file exists
func SessionExists() bool {
	path, err := GetSessionPath()
	if err != nil {
		return false
	}

	_, err = os.Stat(path)
	return err == nil
}
