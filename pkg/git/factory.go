package git

// SessionFactory provides functions to create and manage sessions
// These are implemented by the commands package to avoid circular imports
var (
	newSession      func(string, ...SessionOption) (Session, error)
	loadSession     func(string) (Session, error)
	validateSession func(string) error
	getSessionInfo  func(string) (*SessionConfig, error)
)

// RegisterSessionFactory registers the session factory functions
// This is called by the commands package during initialization
func RegisterSessionFactory(
	newFunc func(string, ...SessionOption) (Session, error),
	loadFunc func(string) (Session, error),
	validateFunc func(string) error,
	infoFunc func(string) (*SessionConfig, error),
) {
	newSession = newFunc
	loadSession = loadFunc
	validateSession = validateFunc
	getSessionInfo = infoFunc
}

// NewSession creates a new Git session with persistent configuration
func NewSession(sessionPath string, opts ...SessionOption) (Session, error) {
	if newSession == nil {
		panic("session factory not registered - import github.com/instruqt/git-exec/pkg/git/commands")
	}
	return newSession(sessionPath, opts...)
}

// LoadSession loads an existing session from a repository path
func LoadSession(sessionPath string) (Session, error) {
	if loadSession == nil {
		panic("session factory not registered - import github.com/instruqt/git-exec/pkg/git/commands")
	}
	return loadSession(sessionPath)
}

// ValidateSession checks if a session at the given path is valid
func ValidateSession(sessionPath string) error {
	if validateSession == nil {
		panic("session factory not registered - import github.com/instruqt/git-exec/pkg/git/commands")
	}
	return validateSession(sessionPath)
}

// GetSessionInfo returns basic information about a session
func GetSessionInfo(sessionPath string) (*SessionConfig, error) {
	if getSessionInfo == nil {
		panic("session factory not registered - import github.com/instruqt/git-exec/pkg/git/commands")
	}
	return getSessionInfo(sessionPath)
}