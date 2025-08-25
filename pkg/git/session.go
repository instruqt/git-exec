package git

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// SessionConfig represents configuration for a Git session
type SessionConfig struct {
	// User information
	UserName  string
	UserEmail string
	
	// Session properties
	WorkingDirectory string
	
	// Metadata (key-value pairs for any use case)
	Metadata map[string]string
}

// Session represents a Git session with persistent configuration
type Session interface {
	Git
	
	// Session-specific methods
	GetSessionConfig() *SessionConfig
	UpdateUser(name, email string) error
	IsValid() bool
	InitRepository() error
	Destroy() error
}

// SessionOption is a functional option for configuring sessions
type SessionOption func(*SessionConfig)

// WithSessionUser sets the user for the session
func SessionWithUser(name, email string) SessionOption {
	return func(c *SessionConfig) {
		c.UserName = name
		c.UserEmail = email
	}
}


// SessionWithMetadata adds custom metadata to the session with a section
func SessionWithMetadata(section, key, value string) SessionOption {
	return func(c *SessionConfig) {
		if c.Metadata == nil {
			c.Metadata = make(map[string]string)
		}
		sectionedKey := fmt.Sprintf("%s.%s", section, key)
		c.Metadata[sectionedKey] = value
	}
}

// SessionWithWorkingDirectory sets the working directory for the session
func SessionWithWorkingDirectory(dir string) SessionOption {
	return func(c *SessionConfig) {
		c.WorkingDirectory = dir
	}
}

// sessionImpl extends gitImpl with session management capabilities
type sessionImpl struct {
	*gitImpl
	config *SessionConfig
}

// NewSession creates a new Git session with persistent configuration
func NewSession(sessionPath string, opts ...SessionOption) (Session, error) {
	// Create base git instance
	g, err := NewGit()
	if err != nil {
		return nil, fmt.Errorf("failed to create git instance: %w", err)
	}
	
	// Initialize session config
	config := &SessionConfig{
		WorkingDirectory: sessionPath,
		Metadata:         make(map[string]string),
	}
	
	// Apply options
	for _, opt := range opts {
		opt(config)
	}
	
	// Create session
	s := &sessionImpl{
		gitImpl: g,
		config:  config,
	}
	
	// Set working directory
	s.SetWorkingDirectory(sessionPath)
	
	// Check if this is an existing git repository
	isExistingRepo := false
	if _, err := os.Stat(sessionPath); err == nil {
		// Directory exists, check if it's a git repository
		cmd := s.newCommand("rev-parse", "--git-dir")
		if _, err := cmd.Execute(); err == nil {
			isExistingRepo = true
		}
	}
	
	if isExistingRepo {
		// Load existing repository configuration
		if err := s.loadConfig(); err != nil {
			// If loading fails, that's ok - just means no session config exists yet
		}
		
		// Update with any new session configuration
		if err := s.persistConfig(); err != nil {
			return nil, fmt.Errorf("failed to update session configuration: %w", err)
		}
	} else {
		// No existing repository - create one and set up session config
		if err := s.ensureRepository(); err != nil {
			return nil, fmt.Errorf("failed to initialize repository: %w", err)
		}
		
		// Write session configuration to .git/config
		if err := s.persistConfig(); err != nil {
			return nil, fmt.Errorf("failed to persist session configuration: %w", err)
		}
	}
	
	return s, nil
}

// LoadSession loads an existing session from a repository path
func LoadSession(sessionPath string) (Session, error) {
	// Check if path exists
	if _, err := os.Stat(sessionPath); os.IsNotExist(err) {
		return nil, fmt.Errorf("session path does not exist: %s", sessionPath)
	}
	
	// Create base git instance
	g, err := NewGit()
	if err != nil {
		return nil, fmt.Errorf("failed to create git instance: %w", err)
	}
	
	// Create session
	s := &sessionImpl{
		gitImpl: g,
		config: &SessionConfig{
			WorkingDirectory: sessionPath,
			Metadata:         make(map[string]string),
		},
	}
	
	// Set working directory
	s.SetWorkingDirectory(sessionPath)
	
	// Verify it's a git repository
	cmd := s.newCommand("rev-parse", "--git-dir")
	if _, err := cmd.Execute(); err != nil {
		return nil, fmt.Errorf("not a git repository: %s", sessionPath)
	}
	
	// Load configuration from .git/config
	if err := s.loadConfig(); err != nil {
		return nil, fmt.Errorf("failed to load configuration: %w", err)
	}
	
	return s, nil
}

// ValidateSession checks if a session at the given path is valid
func ValidateSession(sessionPath string) error {
	s, err := LoadSession(sessionPath)
	if err != nil {
		return fmt.Errorf("failed to load session: %w", err)
	}
	
	if !s.IsValid() {
		return fmt.Errorf("session is not valid")
	}
	
	config := s.GetSessionConfig()
	if config.UserName == "" {
		return fmt.Errorf("user name is missing")
	}
	
	return nil
}

// GetSessionInfo returns basic information about a session
func GetSessionInfo(sessionPath string) (*SessionConfig, error) {
	s, err := LoadSession(sessionPath)
	if err != nil {
		return nil, fmt.Errorf("failed to load session: %w", err)
	}
	
	return s.GetSessionConfig(), nil
}

// GetSessionConfig returns the session configuration
func (s *sessionImpl) GetSessionConfig() *SessionConfig {
	return s.config
}

// UpdateUser updates the user information for the session
func (s *sessionImpl) UpdateUser(name, email string) error {
	s.config.UserName = name
	s.config.UserEmail = email
	return s.persistConfig()
}


// IsValid checks if the session is still valid
func (s *sessionImpl) IsValid() bool {
	// Check if working directory exists
	cmd := s.newCommand("rev-parse", "--git-dir")
	_, err := cmd.Execute()
	return err == nil
}

// InitRepository initializes a repository in the session directory
func (s *sessionImpl) InitRepository() error {
	if err := s.ensureRepository(); err != nil {
		return fmt.Errorf("failed to initialize repository: %w", err)
	}
	
	// Write configuration after repository is initialized
	if err := s.persistConfig(); err != nil {
		return fmt.Errorf("failed to persist configuration: %w", err)
	}
	
	return nil
}

// Destroy removes session-specific configuration
func (s *sessionImpl) Destroy() error {
	// Remove all metadata keys
	for key := range s.config.Metadata {
		cmd := s.newCommand("config", "--local", "--unset", key)
		_, _ = cmd.Execute() // Ignore errors if key doesn't exist
	}
	
	return nil
}

// ensureRepository ensures a git repository exists at the session path
func (s *sessionImpl) ensureRepository() error {
	// Ensure the directory exists
	if s.config.WorkingDirectory != "" {
		if err := os.MkdirAll(s.config.WorkingDirectory, 0755); err != nil {
			return fmt.Errorf("failed to create directory: %w", err)
		}
	}
	
	// Check if .git directory exists
	cmd := s.newCommand("rev-parse", "--git-dir")
	if _, err := cmd.Execute(); err != nil {
		// Repository doesn't exist, initialize it
		initCmd := s.newCommand("init")
		if _, err := initCmd.Execute(); err != nil {
			return fmt.Errorf("failed to initialize repository: %w", err)
		}
	}
	return nil
}

// persistConfig writes session configuration to .git/config
func (s *sessionImpl) persistConfig() error {
	// Set user configuration
	if s.config.UserName != "" {
		if err := s.setConfigValue("user.name", s.config.UserName); err != nil {
			return fmt.Errorf("failed to set user.name: %w", err)
		}
	}
	
	if s.config.UserEmail != "" {
		if err := s.setConfigValue("user.email", s.config.UserEmail); err != nil {
			return fmt.Errorf("failed to set user.email: %w", err)
		}
	}
	
	// Set metadata (keys already include section like "user.id")
	for key, value := range s.config.Metadata {
		if err := s.setConfigValue(key, value); err != nil {
			return fmt.Errorf("failed to set %s: %w", key, err)
		}
	}
	
	return nil
}

// loadConfig reads session configuration from .git/config
func (s *sessionImpl) loadConfig() error {
	// Load user configuration
	if name, err := s.getConfigValue("user.name"); err == nil {
		s.config.UserName = name
	}
	
	if email, err := s.getConfigValue("user.email"); err == nil {
		s.config.UserEmail = email
	}
	
	
	// Load session metadata (look for section.key patterns like user.id, project.name)
	cmd := s.newCommand("config", "--get-regexp", "^[a-zA-Z]+\\.[a-zA-Z]+$")
	output, err := cmd.Execute()
	if err == nil && len(output) > 0 {
		lines := strings.Split(string(output), "\n")
		for _, line := range lines {
			if line == "" {
				continue
			}
			parts := strings.SplitN(line, " ", 2)
			if len(parts) == 2 {
				// Skip user.name and user.email as they're handled separately
				if parts[0] != "user.name" && parts[0] != "user.email" {
					s.config.Metadata[parts[0]] = parts[1]
				}
			}
		}
	}
	
	return nil
}

// setConfigValue sets a git config value
func (s *sessionImpl) setConfigValue(key, value string) error {
	cmd := s.newCommand("config", "--local", key, value)
	_, err := cmd.Execute()
	return err
}

// getConfigValue gets a git config value
func (s *sessionImpl) getConfigValue(key string) (string, error) {
	cmd := s.newCommand("config", "--local", "--get", key)
	output, err := cmd.Execute()
	if err != nil {
		return "", err
	}
	return strings.TrimSpace(string(output)), nil
}

// removeConfig removes a configuration section
func (s *sessionImpl) removeConfig(section string) error {
	cmd := s.newCommand("config", "--local", "--remove-section", section)
	_, err := cmd.Execute()
	// Ignore error if section doesn't exist
	if err != nil && strings.Contains(err.Error(), "No such section") {
		return nil
	}
	return err
}

// Override key git operations to ensure user context is applied

// Commit creates a commit with automatic user attribution
func (s *sessionImpl) Commit(message string, opts ...Option) error {
	// Ensure user context is set
	allOpts := make([]Option, 0, len(opts)+1)
	if s.config.UserName != "" && s.config.UserEmail != "" {
		allOpts = append(allOpts, WithUser(s.config.UserName, s.config.UserEmail))
	}
	allOpts = append(allOpts, opts...)
	
	return s.gitImpl.Commit(message, allOpts...)
}

// Clone implements Clone for sessions, automatically applying session user context
func (s *sessionImpl) Clone(url, destination string, opts ...Option) error {
	// Apply user context from session if available
	allOpts := make([]Option, 0, len(opts)+2)
	if s.config.UserName != "" && s.config.UserEmail != "" {
		allOpts = append(allOpts, 
			WithConfig("user.name", s.config.UserName),
			WithConfig("user.email", s.config.UserEmail),
		)
	}
	allOpts = append(allOpts, opts...)
	
	// Ensure destination directory exists
	if err := os.MkdirAll(filepath.Dir(destination), 0755); err != nil {
		return fmt.Errorf("failed to create destination directory: %w", err)
	}
	
	// Perform the clone
	if err := s.gitImpl.Clone(url, destination, allOpts...); err != nil {
		return err
	}
	
	// If cloning into the session directory, persist session configuration
	if destination == s.config.WorkingDirectory {
		if err := s.persistConfig(); err != nil {
			return fmt.Errorf("failed to persist session configuration: %w", err)
		}
	}
	
	return nil
}