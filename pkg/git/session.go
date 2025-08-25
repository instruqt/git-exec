package git

import (
	"time"
)

// SessionConfig represents configuration for a Git session
type SessionConfig struct {
	// User information
	UserName  string
	UserEmail string
	
	// Instruqt metadata
	UserID    string
	SessionID string
	Created   time.Time
	
	// Session properties
	WorkingDirectory string
	
	// Additional metadata (key-value pairs)
	Metadata map[string]string
}

// Session represents a Git session with persistent configuration
type Session interface {
	Git
	
	// Session-specific methods
	GetConfig() *SessionConfig
	UpdateUser(name, email string) error
	GetSessionID() string
	GetUserID() string
	IsValid() bool
	InitRepository() error
	Destroy() error
}

// SessionOption is a functional option for configuring sessions
type SessionOption func(*SessionConfig)

// WithUser sets the user for the session
func WithUser(name, email string) SessionOption {
	return func(c *SessionConfig) {
		c.UserName = name
		c.UserEmail = email
	}
}

// WithInstruqtMetadata sets Instruqt-specific metadata
func WithInstruqtMetadata(userID, sessionID string, created time.Time) SessionOption {
	return func(c *SessionConfig) {
		c.UserID = userID
		c.SessionID = sessionID
		c.Created = created
	}
}

// WithMetadata adds custom metadata to the session
func WithMetadata(key, value string) SessionOption {
	return func(c *SessionConfig) {
		if c.Metadata == nil {
			c.Metadata = make(map[string]string)
		}
		c.Metadata[key] = value
	}
}

// WithWorkingDirectory sets the working directory for the session
func WithWorkingDirectory(dir string) SessionOption {
	return func(c *SessionConfig) {
		c.WorkingDirectory = dir
	}
}

