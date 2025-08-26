package git

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
func (g *gitImpl) newCommand(operation string, args ...string) Command {
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
func (c *command) ApplyOptions(opts ...Option) {
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
func WithTimeout(timeout time.Duration) Option {
	return func(c Command) {
		c.SetTimeout(timeout)
	}
}

// WithEnv sets an environment variable for the command
func WithEnv(key, value string) Option {
	return func(c Command) {
		c.SetEnv(key, value)
	}
}

// WithWorkingDirectory sets the working directory for the command
func WithWorkingDirectory(dir string) Option {
	return func(c Command) {
		c.SetWorkingDir(dir)
	}
}

// WithStdin sets the stdin for the command
func WithStdin(input string) Option {
	return func(c Command) {
		c.SetStdin(input)
	}
}

// WithAuth sets authentication for remote operations
// This is generic enough to be used by multiple commands (clone, fetch, push, pull)
func WithAuth(token string) Option {
	return func(c Command) {
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
func WithUser(name, email string) Option {
	return func(c Command) {
		c.SetEnv("GIT_AUTHOR_NAME", name)
		c.SetEnv("GIT_AUTHOR_EMAIL", email)
		c.SetEnv("GIT_COMMITTER_NAME", name)
		c.SetEnv("GIT_COMMITTER_EMAIL", email)
	}
}

// WithQuiet adds the -q/--quiet flag (common across many commands)
func WithQuiet() Option {
	return func(c Command) {
		c.AddArgs("--quiet")
	}
}

// WithVerbose adds the -v/--verbose flag (common across many commands)
func WithVerbose() Option {
	return func(c Command) {
		c.AddArgs("--verbose")
	}
}

// WithArgs adds arbitrary arguments to the command (escape hatch for unsupported options)
func WithArgs(args ...string) Option {
	return func(c Command) {
		c.AddArgs(args...)
	}
}

// WithConfig sets a git config value for the command
// This adds -c key=value to the git command
func WithConfig(key, value string) Option {
	return func(c Command) {
		// Git config options must come before the subcommand
		args := c.GetArgs()
		if len(args) > 0 {
			// Insert config option before the subcommand
			newArgs := []string{"-c", fmt.Sprintf("%s=%s", key, value)}
			newArgs = append(newArgs, args...)
			c.SetArgs(newArgs)
		}
	}
}

// WithConfigs sets multiple git config values for the command
func WithConfigs(configs map[string]string) Option {
	return func(c Command) {
		args := c.GetArgs()
		if len(args) > 0 {
			// Build all config options
			var configArgs []string
			for key, value := range configs {
				configArgs = append(configArgs, "-c", fmt.Sprintf("%s=%s", key, value))
			}
			// Insert config options before the subcommand
			configArgs = append(configArgs, args...)
			c.SetArgs(configArgs)
		}
	}
}

// Add-specific options

// AddWithForce allows adding ignored files
func AddWithForce() Option {
	return WithArgs("--force")
}

// AddWithDryRun shows what would be added without actually adding
func AddWithDryRun() Option {
	return WithArgs("--dry-run")
}

// AddWithVerbose shows files as they are added
func AddWithVerbose() Option {
	return WithArgs("--verbose")
}

// AddWithAll stages all changes (modifications, deletions, new files)
func AddWithAll() Option {
	return WithArgs("--all")
}

// AddWithUpdate stages modifications and deletions, but not new files
func AddWithUpdate() Option {
	return WithArgs("--update")
}

// AddWithNoIgnoreRemoval doesn't ignore removed files
func AddWithNoIgnoreRemoval() Option {
	return WithArgs("--no-ignore-removal")
}

// AddWithIgnoreErrors continues adding files even if some fail
func AddWithIgnoreErrors() Option {
	return WithArgs("--ignore-errors")
}

// AddWithIntent records only the fact that a path will be added later
func AddWithIntent() Option {
	return WithArgs("--intent-to-add")
}

// AddWithPatch interactively choose hunks to add
func AddWithPatch() Option {
	return WithArgs("--patch")
}

// Status-specific options

// StatusWithShort gives output in short format
func StatusWithShort() Option {
	return WithArgs("--short")
}

// StatusWithBranch shows branch information
func StatusWithBranch() Option {
	return WithArgs("--branch")
}

// StatusWithPorcelain gives porcelain output (default for this implementation)
func StatusWithPorcelain() Option {
	return WithArgs("--porcelain")
}

// StatusWithLong gives output in long format (default Git behavior)
func StatusWithLong() Option {
	return WithArgs("--long")
}

// StatusWithShowStash shows stash information
func StatusWithShowStash() Option {
	return WithArgs("--show-stash")
}

// StatusWithAheadBehind shows ahead/behind counts
func StatusWithAheadBehind() Option {
	return WithArgs("--ahead-behind")
}

// StatusWithUntrackedFiles controls how untracked files are shown
func StatusWithUntrackedFiles(mode string) Option {
	return WithArgs("--untracked-files=" + mode)
}

// StatusWithIgnoredFiles shows ignored files
func StatusWithIgnoredFiles() Option {
	return WithArgs("--ignored")
}

// Commit-specific options

// CommitWithAuthor sets the author for the commit (alias for WithUser)
func CommitWithAuthor(name, email string) Option {
	return WithUser(name, email)
}

// CommitWithAll automatically stages all modified and deleted files
func CommitWithAll() Option {
	return WithArgs("--all")
}

// CommitWithAmend replaces the tip of the current branch
func CommitWithAmend() Option {
	return WithArgs("--amend")
}

// CommitWithNoEdit uses the previous commit message without launching an editor
func CommitWithNoEdit() Option {
	return WithArgs("--no-edit")
}

// CommitWithAllowEmpty allows creating a commit with no changes
func CommitWithAllowEmpty() Option {
	return WithArgs("--allow-empty")
}

// CommitWithAllowEmptyMessage allows a commit with an empty message
func CommitWithAllowEmptyMessage() Option {
	return WithArgs("--allow-empty-message")
}

// CommitWithSignoff adds a Signed-off-by line
func CommitWithSignoff() Option {
	return WithArgs("--signoff")
}

// CommitWithGPGSign signs the commit with GPG
func CommitWithGPGSign(keyid string) Option {
	if keyid == "" {
		return WithArgs("--gpg-sign")
	}
	return WithArgs("--gpg-sign=" + keyid)
}

// CommitWithNoVerify bypasses pre-commit and commit-msg hooks
func CommitWithNoVerify() Option {
	return WithArgs("--no-verify")
}

// Log-specific options

// LogWithMaxCount limits the number of commits to show
func LogWithMaxCount(count string) Option {
	return WithArgs("--max-count", count)
}

// LogWithOneline shows commits in oneline format
func LogWithOneline() Option {
	return WithArgs("--oneline")
}

// LogWithGraph shows a text-based graphical representation
func LogWithGraph() Option {
	return WithArgs("--graph")
}

// LogWithStat shows diffstat for each commit
func LogWithStat() Option {
	return WithArgs("--stat")
}

// Checkout-specific options

// CheckoutWithBranch specifies the branch to checkout
func CheckoutWithBranch(branch string) Option {
	return WithArgs(branch)
}

// CheckoutWithCreate creates a new branch and checks it out
func CheckoutWithCreate(branch string) Option {
	return WithArgs("-b", branch)
}

// CheckoutWithCreateFrom creates a new branch from a specific commit and checks it out
func CheckoutWithCreateFrom(branch, startPoint string) Option {
	return WithArgs("-b", branch, startPoint)
}

// CheckoutWithForce forces the checkout (discards local changes)
func CheckoutWithForce() Option {
	return WithArgs("--force")
}

// CheckoutWithCommit checks out the specified commit
func CheckoutWithCommit(commit string) Option {
	return WithArgs(commit)
}


// CheckoutWithOrphan creates an orphan branch
func CheckoutWithOrphan(branch string) Option {
	return WithArgs("--orphan", branch)
}

// Merge-specific options

// MergeWithBranch specifies the branch to merge
func MergeWithBranch(branch string) Option {
	return WithArgs(branch)
}

// MergeWithNoFF creates a merge commit even when a fast-forward is possible
func MergeWithNoFF() Option {
	return WithArgs("--no-ff")
}

// MergeWithFFOnly aborts unless a fast-forward is possible
func MergeWithFFOnly() Option {
	return WithArgs("--ff-only")
}

// MergeWithSquash creates a single commit instead of a merge commit
func MergeWithSquash() Option {
	return WithArgs("--squash")
}

// MergeWithStrategy specifies the merge strategy
func MergeWithStrategy(strategy string) Option {
	return WithArgs("-s", strategy)
}

// MergeWithCommitMessage specifies a custom merge commit message
func MergeWithCommitMessage(message string) Option {
	return WithArgs("-m", message)
}

// Init-specific options

// InitWithBare creates a bare repository
func InitWithBare() Option {
	return WithArgs("--bare")
}

// InitWithTemplate specifies a template directory
func InitWithTemplate(templateDir string) Option {
	return WithArgs("--template", templateDir)
}

// InitWithSeparateGitDir creates the .git directory at a separate location
func InitWithSeparateGitDir(gitDir string) Option {
	return WithArgs("--separate-git-dir", gitDir)
}

// InitWithSharedRepo sets up a shared repository
func InitWithSharedRepo(permissions string) Option {
	if permissions == "" {
		return WithArgs("--shared")
	}
	return WithArgs("--shared=" + permissions)
}

// Clone-specific options

// CloneWithBare creates a bare clone of the repository
func CloneWithBare() Option {
	return WithArgs("--bare")
}

// CloneWithDepth creates a shallow clone with specified depth
func CloneWithDepth(depth int) Option {
	return WithArgs("--depth", fmt.Sprintf("%d", depth))
}

// CloneWithBranch clones only a specific branch
func CloneWithBranch(branch string) Option {
	return WithArgs("--branch", branch)
}

// CloneWithSingleBranch clones only a single branch
func CloneWithSingleBranch() Option {
	return WithArgs("--single-branch")
}

// Config-specific options

// ConfigWithLocalScope operates on repository-specific config
func ConfigWithLocalScope() Option {
	return WithArgs("--local")
}

// ConfigWithGlobalScope operates on user-specific config
func ConfigWithGlobalScope() Option {
	return WithArgs("--global")
}

// ConfigWithSystemScope operates on system-wide config
func ConfigWithSystemScope() Option {
	return WithArgs("--system")
}

// ConfigWithAllScopes lists config from all scopes (for ListConfig)
func ConfigWithAllScopes() Option {
	return WithArgs("--show-scope")
}

// ConfigWithShowOrigin shows the origin file for each config (for ListConfig)
func ConfigWithShowOrigin() Option {
	return WithArgs("--show-origin")
}

