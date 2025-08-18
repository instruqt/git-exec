package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	gitpkg "github.com/instruqt/git-exec/pkg/git"
	"github.com/instruqt/git-exec/pkg/git/errors"
)

// command represents a git command to be executed
type command struct {
	gitPath     string
	args        []string
	workingDir  string
	env         map[string]string
	timeout     time.Duration
	stdin       *bytes.Buffer
	noStderr    bool // Some commands write normal output to stderr (e.g., fetch)
}

// newCommand creates a new command with the given git operation
func (g *git) newCommand(operation string, args ...string) gitpkg.Command {
	cmd := &command{
		gitPath:    g.path,
		args:       append([]string{operation}, args...),
		workingDir: g.wd,
		env:        make(map[string]string),
		timeout:    2 * time.Minute, // Default timeout
	}
	return cmd
}

// Execute runs the git command and returns the output
func (c *command) Execute() ([]byte, error) {
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, c.gitPath, c.args...)
	
	if c.workingDir != "" {
		cmd.Dir = c.workingDir
	}

	// Build environment
	if len(c.env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range c.env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Set stdin if provided
	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}

	// Capture both stdout and stderr
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	
	// Some git commands write normal output to stderr (e.g., fetch, push)
	// In those cases, we should return stderr as the output
	if c.noStderr && err == nil {
		return stderr.Bytes(), nil
	}

	if err != nil {
		// Check if it's an exit error and create a GitError
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, &errors.GitError{
				Command:  c.args,
				ExitCode: exitError.ExitCode(),
				Stderr:   stderr.String(),
				Stdout:   stdout.String(),
			}
		}
		return nil, err
	}

	return stdout.Bytes(), nil
}

// ExecuteCombined runs the git command and returns combined stdout and stderr
func (c *command) ExecuteCombined() ([]byte, error) {
	ctx := context.Background()
	if c.timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, c.timeout)
		defer cancel()
	}

	cmd := exec.CommandContext(ctx, c.gitPath, c.args...)
	
	if c.workingDir != "" {
		cmd.Dir = c.workingDir
	}

	// Build environment
	if len(c.env) > 0 {
		cmd.Env = os.Environ()
		for k, v := range c.env {
			cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", k, v))
		}
	}

	// Set stdin if provided
	if c.stdin != nil {
		cmd.Stdin = c.stdin
	}

	output, err := cmd.CombinedOutput()
	
	if err != nil {
		// Check if it's an exit error and create a GitError
		if exitError, ok := err.(*exec.ExitError); ok {
			return nil, &errors.GitError{
				Command:  c.args,
				ExitCode: exitError.ExitCode(),
				Stderr:   string(output),
			}
		}
		return nil, err
	}

	return output, nil
}

// ApplyOptions applies all options to the command
func (c *command) ApplyOptions(opts ...gitpkg.Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// String returns the command as it would be executed
func (c *command) String() string {
	return fmt.Sprintf("%s %s", c.gitPath, strings.Join(c.args, " "))
}

// Interface method implementations for option configuration
func (c *command) SetTimeout(timeout time.Duration) {
	c.timeout = timeout
}

func (c *command) SetEnv(key, value string) {
	c.env[key] = value
}

func (c *command) SetWorkingDir(dir string) {
	c.workingDir = dir
}

func (c *command) SetStdin(input string) {
	c.stdin = bytes.NewBufferString(input)
}

func (c *command) AddArgs(args ...string) {
	c.args = append(c.args, args...)
}

func (c *command) GetArgs() []string {
	return c.args
}

func (c *command) SetArgs(args []string) {
	c.args = args
}

func (c *command) ExecuteWithStderr() ([]byte, error) {
	c.noStderr = true
	return c.Execute()
}

// Generic Options that apply to all commands

// WithTimeout sets a custom timeout for the command
func WithTimeout(timeout time.Duration) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.SetTimeout(timeout)
	}
}

// WithEnv sets an environment variable for the command
func WithEnv(key, value string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.SetEnv(key, value)
	}
}

// WithWorkingDirectory sets the working directory for the command
func WithWorkingDirectory(dir string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.SetWorkingDir(dir)
	}
}

// WithStdin sets the stdin for the command
func WithStdin(input string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.SetStdin(input)
	}
}

// WithAuth sets authentication for remote operations
// This is generic enough to be used by multiple commands (clone, fetch, push, pull)
func WithAuth(token string) gitpkg.Option {
	return func(c gitpkg.Command) {
		// Use the token as a bearer token in the Authorization header
		// This is typically done via GIT_ASKPASS or credential helpers
		c.SetEnv("GIT_ASKPASS", "echo")
		c.SetEnv("GIT_TERMINAL_PROMPT", "0")
		
		// For GitHub token auth, we can use the token directly
		// This would need to be handled differently based on the remote URL
		// For now, we'll set it as an environment variable that could be used
		// by a credential helper
		c.SetEnv("GITHUB_TOKEN", token)
	}
}

// WithUser sets the user for commits (used by multiple commands like commit, merge, etc.)
func WithUser(name, email string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.SetEnv("GIT_AUTHOR_NAME", name)
		c.SetEnv("GIT_AUTHOR_EMAIL", email)
		c.SetEnv("GIT_COMMITTER_NAME", name)
		c.SetEnv("GIT_COMMITTER_EMAIL", email)
	}
}

// WithQuiet adds the -q/--quiet flag (common across many commands)
func WithQuiet() gitpkg.Option {
	return func(c gitpkg.Command) {
		c.AddArgs("--quiet")
	}
}

// WithVerbose adds the -v/--verbose flag (common across many commands)
func WithVerbose() gitpkg.Option {
	return func(c gitpkg.Command) {
		c.AddArgs("--verbose")
	}
}

// WithArgs adds arbitrary arguments to the command (escape hatch for unsupported options)
func WithArgs(args ...string) gitpkg.Option {
	return func(c gitpkg.Command) {
		c.AddArgs(args...)
	}
}