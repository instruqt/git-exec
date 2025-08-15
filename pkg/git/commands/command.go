package commands

import (
	"bytes"
	"context"
	"fmt"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/instruqt/git-exec/pkg/git/errors"
)

// Command represents a git command to be executed
type Command struct {
	gitPath     string
	args        []string
	workingDir  string
	env         map[string]string
	timeout     time.Duration
	stdin       *bytes.Buffer
	noStderr    bool // Some commands write normal output to stderr (e.g., fetch)
}

// Option is a functional option for configuring git commands
type Option func(*Command)

// newCommand creates a new command with the given git operation
func (g *git) newCommand(operation string, args ...string) *Command {
	cmd := &Command{
		gitPath:    g.path,
		args:       append([]string{operation}, args...),
		workingDir: g.wd,
		env:        make(map[string]string),
		timeout:    2 * time.Minute, // Default timeout
	}
	return cmd
}

// Execute runs the git command and returns the output
func (c *Command) Execute() ([]byte, error) {
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
func (c *Command) ExecuteCombined() ([]byte, error) {
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

// Apply applies all options to the command
func (c *Command) applyOptions(opts ...Option) {
	for _, opt := range opts {
		opt(c)
	}
}

// Helper function to parse git command output that may be on stderr
func (c *Command) executeWithStderr() ([]byte, error) {
	c.noStderr = true
	return c.Execute()
}

// String returns the command as it would be executed
func (c *Command) String() string {
	return fmt.Sprintf("%s %s", c.gitPath, strings.Join(c.args, " "))
}

// Generic Options that apply to all commands

// WithTimeout sets a custom timeout for the command
func WithTimeout(timeout time.Duration) Option {
	return func(c *Command) {
		c.timeout = timeout
	}
}

// WithEnv sets an environment variable for the command
func WithEnv(key, value string) Option {
	return func(c *Command) {
		c.env[key] = value
	}
}

// WithWorkingDirectory sets the working directory for the command
func WithWorkingDirectory(dir string) Option {
	return func(c *Command) {
		c.workingDir = dir
	}
}

// WithStdin sets the stdin for the command
func WithStdin(input string) Option {
	return func(c *Command) {
		c.stdin = bytes.NewBufferString(input)
	}
}

// WithAuth sets authentication for remote operations
// This is generic enough to be used by multiple commands (clone, fetch, push, pull)
func WithAuth(token string) Option {
	return func(c *Command) {
		// Use the token as a bearer token in the Authorization header
		// This is typically done via GIT_ASKPASS or credential helpers
		c.env["GIT_ASKPASS"] = "echo"
		c.env["GIT_TERMINAL_PROMPT"] = "0"
		
		// For GitHub token auth, we can use the token directly
		// This would need to be handled differently based on the remote URL
		// For now, we'll set it as an environment variable that could be used
		// by a credential helper
		c.env["GITHUB_TOKEN"] = token
	}
}

// WithUser sets the user for commits (used by multiple commands like commit, merge, etc.)
func WithUser(name, email string) Option {
	return func(c *Command) {
		c.env["GIT_AUTHOR_NAME"] = name
		c.env["GIT_AUTHOR_EMAIL"] = email
		c.env["GIT_COMMITTER_NAME"] = name
		c.env["GIT_COMMITTER_EMAIL"] = email
	}
}

// WithQuiet adds the -q/--quiet flag (common across many commands)
func WithQuiet() Option {
	return func(c *Command) {
		c.args = append(c.args, "--quiet")
	}
}

// WithVerbose adds the -v/--verbose flag (common across many commands)
func WithVerbose() Option {
	return func(c *Command) {
		c.args = append(c.args, "--verbose")
	}
}

// WithArgs adds arbitrary arguments to the command (escape hatch for unsupported options)
func WithArgs(args ...string) Option {
	return func(c *Command) {
		c.args = append(c.args, args...)
	}
}