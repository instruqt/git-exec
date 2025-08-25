package git

import "time"

// Command interface defines the contract for git command execution
type Command interface {
	Execute() ([]byte, error)
	ExecuteCombined() ([]byte, error)
	ExecuteWithStderr() ([]byte, error)
	ApplyOptions(opts ...Option)
	// Internal methods for option configuration
	SetTimeout(timeout time.Duration)
	SetEnv(key, value string)
	SetWorkingDir(dir string)
	SetStdin(input string)
	AddArgs(args ...string)
	// Internal access methods
	GetArgs() []string
	SetArgs(args []string)
}

// Option is a functional option for configuring git commands
type Option func(Command)