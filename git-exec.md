# Git-Exec Library Documentation

## Overview

Git-exec is a Go library that provides a clean, type-safe wrapper around the Git CLI. Designed specifically for the Instruqt VCS service, it uses a functional options pattern to offer both simplicity for basic operations and extensibility for complex scenarios.

## Current Implementation Status

The git-exec library has a **modern, production-ready architecture** with the following features:

### ✅ Completed Features

**Architecture & Infrastructure**:
- Modern package structure (`pkg/git/`, `pkg/git/commands/`, `pkg/git/types/`, `pkg/git/errors/`)
- Type-safe functional options pattern with `Option func(Command)`
- Reusable command infrastructure with centralized execution
- Structured error handling with `GitError` types
- Environment variable support for authentication
- Timeout and context support for all operations

**Implemented Commands (30+)**:
- **Core Operations**: Clone, Init, Add, Commit, Status, Reset
- **Remote Management**: AddRemote, RemoveRemote, SetRemoteURL, ListRemotes  
- **Branch Operations**: ListBranches, CreateBranch, DeleteBranch, SetUpstream, Checkout
- **Repository Inspection**: Log, Show, Diff (with sophisticated parsing)
- **Synchronization**: Fetch, Pull, Push
- **Tag Operations**: Tag, ListTags, DeleteTag, PushTags, DeleteRemoteTag
- **Advanced Operations**: Merge, Rebase, Revert, Reflog, Config, Remove

### 🚧 Commands with Basic Implementation

The following commands have basic structure but need completion:
- `Merge()` - Has extensive options but needs conflict parsing
- `Checkout()` - Basic implementation, needs output parsing
- `Config()` - Basic implementation, needs get/set logic refinement
- `Revert()`, `Reflog()`, `Remove()` - Have option definitions but basic execution

### ✅ Recently Implemented

- **Session Management**: Complete persistent user sessions with `.git/config` storage
  - Smart repository detection (auto-detects existing repos vs new directories)
  - Automatic user attribution for all session operations
  - Custom metadata storage in `[session]` section
  - Session validation, loading, and destruction
- **Enhanced Options Pattern**: 
  - `WithConfig()` for temporary git config on any operation
  - `WithConfigs()` for multiple config values at once
  - Clean, consistent API: `NewSession()` instead of `NewSessionWithConfig()`

### ❌ Not Yet Implemented

- **Conflict Resolution**: Parsing and handling merge conflicts  
- **Bare Repository**: Specific handling for bare repository operations

## API Design

### Core Interface

The library provides a clean `Git` interface with all operations:

```go
type Git interface {
    // Working directory management
    SetWorkingDirectory(wd string)
    
    // Repository operations
    Init(path string, options ...Option) error
    Clone(url, destination string, options ...Option) error
    
    // File operations
    Add(files []string, options ...Option) error
    Reset(files []string, options ...Option) error
    Remove(options ...Option) error
    
    // Commit operations
    Commit(message string, options ...Option) error
    Revert(options ...Option) error
    
    // Branch operations
    ListBranches(options ...Option) ([]types.Branch, error)
    CreateBranch(branch string, options ...Option) error
    DeleteBranch(branch string, options ...Option) error
    SetUpstream(branch string, remote string, options ...Option) error
    Checkout(options ...Option) error
    
    // Remote operations
    AddRemote(name, url string, options ...Option) error
    RemoveRemote(name string, options ...Option) error
    SetRemoteURL(name, url string, options ...Option) error
    ListRemotes(options ...Option) ([]types.Remote, error)
    
    // Synchronization
    Fetch(options ...Option) ([]types.Remote, error)
    Pull(options ...Option) (*types.MergeResult, error)
    Push(options ...Option) ([]types.Remote, error)
    
    // Tag operations
    Tag(name string, options ...Option) error
    ListTags(options ...Option) ([]string, error)
    DeleteTag(name string, options ...Option) error
    PushTags(remote string, options ...Option) ([]types.Remote, error)
    DeleteRemoteTag(remote, tagName string, options ...Option) error
    
    // Inspection
    Status(options ...Option) ([]types.File, error)
    Diff(options ...Option) ([]types.Diff, error)
    Show(object string, options ...Option) (*types.Log, error)
    Log(options ...Option) ([]types.Log, error)
    
    // Merge and rebase
    Merge(options ...Option) error
    Rebase(options ...Option) error
    
    // Configuration
    Config(key string, value string, options ...Option) error
    Reflog(options ...Option) error
}
```

### Options Pattern

The library uses functional options for all commands, providing type safety and flexibility:

```go
// Session operations (persistent user context)
session := git.NewSession("/path/to/project",
    git.WithUser("Jane Developer", "jane@instruqt.com"),
    git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
    git.WithMetadata("track", "git-basics"),
)

err := session.Add([]string{"file.txt"})
err := session.Commit("Fix authentication bug") // Uses Jane's context automatically

// One-off operations (temporary config)
err := git.Add([]string{"*.go"}, 
    AddWithDryRun(),
    AddWithVerbose(),
    WithTimeout(30*time.Second))

// Clone with temporary authentication
err := git.Clone(repoURL, destPath,
    CloneWithBare(),
    CloneWithBranch("main"),
    CloneWithDepth(1),
    WithAuth(token))

// Commit with temporary user config
err := git.Commit("Quick fix",
    commands.WithConfig("user.name", "Temp Worker"),
    commands.WithConfig("user.email", "temp@example.com"),
    CommitWithSignoff())

// Multiple temporary configs at once
configs := map[string]string{
    "user.name": "Jane Developer",
    "user.email": "jane@example.com",
    "commit.gpgsign": "false",
}
err := git.Commit("Multi-config commit", commands.WithConfigs(configs))
```

### Generic Options

Available for all commands:

```go
WithTimeout(duration)        // Set command timeout
WithAuth(token)              // GitHub token authentication  
WithUser(name, email)        // Git user attribution (environment variables)
WithConfig(key, value)       // Temporary git config (-c key=value)
WithConfigs(map[string]string) // Multiple git config values at once
WithEnv(key, value)          // Environment variables
WithWorkingDirectory(dir)    // Custom working directory
WithStdin(input)             // Command input
WithQuiet()                  // Suppress output
WithVerbose()                // Verbose output
WithArgs(args...)            // Escape hatch for unsupported options
```

### Session Options

Available for session creation:

```go
WithUser(name, email)        // Set session user (persistent in .git/config)
WithInstruqtMetadata(userID, sessionID, created) // Instruqt-specific metadata
WithMetadata(key, value)     // Custom metadata (stored in [session] section)
WithWorkingDirectory(dir)    // Session working directory
```

### Command-Specific Options

Each command has its own set of specific options. Examples:

**Clone Options**:
```go
CloneWithBare()              // Create bare repository
CloneWithBranch(branch)      // Specific branch
CloneWithDepth(depth)        // Shallow clone
CloneWithSingleBranch()      // Single branch only
CloneWithNoCheckout()        // Skip checkout
CloneWithRecurse()           // Include submodules
```

**Commit Options**:
```go
CommitWithAll()              // Commit all changed files
CommitWithAmend()            // Amend previous commit
CommitWithSignoff()          // Add Signed-off-by
CommitWithNoVerify()         // Skip hooks
CommitWithAllowEmpty()       // Allow empty commit
CommitWithAuthor(author)     // Override author
```

**Push Options**:
```go
PushWithForce()              // Force push
PushWithSetUpstream()        // Set upstream branch
PushWithTags()               // Push tags
PushWithDryRun()             // Dry run
PushWithAtomic()             // Atomic push
```

**Merge Options**:
```go
MergeWithBranch(branch)      // Merge specific branch
MergeWithCommit(commit)      // Merge commit SHA
MergeWithNoFF()              // Force merge commit
MergeWithFFOnly()            // Fast-forward only
MergeWithSquash()            // Squash commits
MergeWithStrategy(strategy)  // Use specific strategy
MergeWithMessage(msg)        // Custom merge message
MergeWithAbort()             // Abort ongoing merge
MergeWithContinue()          // Continue after resolving
```

## Merge Operations and Conflict Resolution

The library provides comprehensive merge functionality with automatic conflict detection and resolution helpers.

### Basic Merge Operations

```go
// Simple branch merge
result, err := git.Merge(MergeWithBranch("feature-branch"))
if err != nil {
    log.Fatal(err)
}

if result.Success {
    fmt.Printf("Successfully merged %s\n", result.MergedBranch)
    if result.FastForward {
        fmt.Println("Fast-forward merge")
    } else {
        fmt.Printf("Merge commit created using %s strategy\n", result.Strategy)
    }
} else {
    fmt.Printf("Merge failed: %s\n", result.AbortReason)
}
```

### Conflict Detection and Handling

```go
// Merge with conflict detection
result, err := git.Merge(MergeWithBranch("conflicting-branch"))
if err != nil {
    log.Fatal(err)
}

if !result.Success && len(result.Conflicts) > 0 {
    fmt.Printf("Merge conflicts detected in %d files:\n", len(result.ConflictedFiles))
    
    for _, conflict := range result.Conflicts {
        fmt.Printf("File: %s (Status: %s)\n", conflict.Path, conflict.Status)
        
        for i, section := range conflict.Sections {
            fmt.Printf("  Conflict %d (lines %d-%d):\n", i+1, section.StartLine, section.EndLine)
            fmt.Printf("    Our version:\n%s\n", section.OurContent)
            fmt.Printf("    Their version:\n%s\n", section.TheirContent)
        }
    }
}
```

### Conflict Resolution Strategies

#### Use Ours Strategy
```go
// Resolve all conflicts by keeping our version
resolutions := []types.ConflictResolution{
    {
        FilePath: "conflicted-file.txt",
        UseOurs:  true,
    },
}

err := git.ResolveConflicts(resolutions)
if err != nil {
    log.Fatal(err)
}

// Continue the merge
err = git.MergeContinue()
```

#### Use Theirs Strategy
```go
// Resolve conflicts by keeping their version
resolutions := []types.ConflictResolution{
    {
        FilePath: "conflicted-file.txt",
        UseTheirs: true,
    },
}

err := git.ResolveConflicts(resolutions)
```

#### Custom Resolution
```go
// Provide custom resolution for specific sections
resolutions := []types.ConflictResolution{
    {
        FilePath: "src/main.go",
        Custom:   true,
        Sections: []types.ResolvedSection{
            {
                SectionIndex: 0,
                Resolution:   "// This is our custom resolution\nfunc main() {\n    fmt.Println(\"Hello, World!\")\n}",
            },
        },
    },
}

err := git.ResolveConflicts(resolutions)
```

### Merge Result Information

The `MergeResult` provides comprehensive information about the merge operation:

```go
type MergeResult struct {
    Success          bool                 // Whether merge completed successfully
    FastForward      bool                 // Whether it was a fast-forward merge
    MergeCommit      string               // SHA of merge commit (if created)
    MergedBranch     string               // Name of branch that was merged
    BaseBranch       string               // Base branch for the merge
    Strategy         string               // Merge strategy used (recursive, ort, octopus)
    ConflictedFiles  []string             // List of files with conflicts
    Conflicts        []ConflictFile       // Detailed conflict information
    Stats            MergeStats           // File change statistics
    AbortReason      string               // Reason if merge failed
}

type ConflictFile struct {
    Path      string               // File path
    Status    ConflictStatus       // Type of conflict (both_modified, etc.)
    Sections  []ConflictSection    // Individual conflict sections
    Content   string               // Raw file content with markers
}

type ConflictSection struct {
    StartLine    int      // Line number where conflict starts
    EndLine      int      // Line number where conflict ends
    OurContent   string   // Content from our branch
    TheirContent string   // Content from their branch
    BaseContent  string   // Original content (with diff3 style)
}
```

### Aborting Merges

```go
// Abort an ongoing merge and return to pre-merge state
err := git.MergeAbort()
if err != nil {
    log.Fatal(err)
}

// Verify clean state
files, err := git.Status()
if err != nil {
    log.Fatal(err)
}

// Check that no files are in conflicted state
for _, file := range files {
    if file.Status == "UU" { // UU indicates unmerged conflict
        fmt.Printf("Warning: %s still shows conflicts\n", file.Name)
    }
}
```

### Advanced Merge Options

```go
// Merge with specific strategy and custom message
result, err := git.Merge(
    MergeWithBranch("feature"),
    MergeWithStrategy("recursive"),
    MergeWithStrategyOption("ours"),
    MergeWithMessage("Merge feature branch with conflict resolution"),
    MergeWithNoEdit(),
)

// Force merge commit even on fast-forward
result, err := git.Merge(
    MergeWithBranch("hotfix"),
    MergeWithNoFF(),
    MergeWithMessage("Hotfix: Critical security update"),
)

// Squash merge (combines all commits into single commit)
result, err := git.Merge(
    MergeWithBranch("feature-branch"),
    MergeWithSquash(),
)
```

## Session Metadata System

Sessions store persistent metadata in `.git/config` that survives service restarts and git operations:

### Storage Location

Metadata is stored in three sections of `.git/config`:

```ini
[user]                        # Session user context
    name = Jane Developer
    email = jane@instruqt.com

[instruqt]                    # Instruqt-specific metadata
    userid = user-123
    sessionid = session-456
    created = 2025-08-25T10:00:20+02:00

[session]                     # Custom application metadata
    track = git-basics
    environment = production
    project = web-app
    team = frontend
    priority = high
```

### Metadata Usage

```go
// Create session with metadata
session, err := git.NewSession("/path/to/project",
    git.WithUser("Jane Developer", "jane@instruqt.com"),
    git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
    git.WithMetadata("track", "git-basics"),
    git.WithMetadata("environment", "production"),
    git.WithMetadata("project", "web-app"),
    git.WithMetadata("team", "frontend"),
    git.WithMetadata("priority", "high"),
)

// Retrieve metadata
config := session.GetConfig()
fmt.Printf("Track: %s\n", config.Metadata["track"])        // git-basics
fmt.Printf("Environment: %s\n", config.Metadata["environment"]) // production
fmt.Printf("Project: %s\n", config.Metadata["project"])    // web-app
```

### Common Use Cases

- **Instruqt Tracks**: Store track name, challenge info, user progress
- **Multi-tenant Applications**: Store tenant ID, organization info
- **CI/CD Pipelines**: Store pipeline ID, build number, deployment target
- **Project Management**: Store project type, priority, team assignment
- **Feature Flags**: Store debug settings, experimental features
- **Audit Trails**: Store additional context beyond git's built-in data

### Session vs Temporary Config

| Feature | Session Config | Temporary Config |
|---------|---------------|------------------|
| **Storage** | `.git/config` (persistent) | Command-only (temporary) |
| **Usage** | `git.WithUser()` in session | `commands.WithConfig()` per operation |
| **Survives** | Service restarts, container rebuilds | Single command only |
| **Best for** | Long-running sessions, user attribution | One-off operations, testing |

```go
// Session config - persistent
session := git.NewSession("/path", git.WithUser("Alice", "alice@company.com"))
session.Commit("My commit") // Always uses Alice's info

// Temporary config - one command only  
git.Commit("Quick fix", commands.WithConfig("user.name", "Bob"))
git.Commit("Another commit") // Uses repository default, not Bob
```

## Structured Types

The library uses structured types for Git objects:

```go
// File status
type File struct {
    Path       string
    Status     string
    OldPath    string // For renames
    IsStaged   bool
    IsTracked  bool
}

// Commit information
type Log struct {
    Hash      string
    Author    Author
    Committer Author
    Message   string
    Date      time.Time
    Body      string
    Files     []FileStat
}

// Branch information
type Branch struct {
    Name      string
    Hash      string
    IsCurrent bool
    IsRemote  bool
    Upstream  string
}

// Remote information
type Remote struct {
    Name     string
    FetchURL string
    PushURL  string
    Branches []string
}

// Merge result
type MergeResult struct {
    Success       bool
    FastForward   bool
    MergeCommit   string
    Conflicts     []string
    MergedBranch  string
}

// Diff information
type Diff struct {
    FileA       string
    FileB       string
    Hash        string
    Status      string
    Insertions  int
    Deletions   int
    IsBinary    bool
    Hunks       []Hunk
}
```

## Error Handling

Structured error handling with `GitError`:

```go
type GitError struct {
    Command  []string
    ExitCode int
    Stderr   string
    Stdout   string
}

// Usage
result, err := git.Merge(MergeWithBranch("feature"))
if gitErr, ok := err.(*errors.GitError); ok {
    switch gitErr.ExitCode {
    case 1:
        // Handle merge conflict
        fmt.Println("Merge conflict:", gitErr.Stderr)
    case 128:
        // Handle invalid operation
        fmt.Println("Invalid operation:", gitErr.Stderr)
    }
}
```

## Command Execution Infrastructure

The library has a robust command execution system:

```go
// Internal command structure
type command struct {
    gitPath    string
    args       []string
    workingDir string
    env        map[string]string
    timeout    time.Duration
    stdin      *bytes.Buffer
    noStderr   bool
}

// Execution methods
cmd.Execute()         // Returns stdout, creates GitError on failure
cmd.ExecuteCombined() // Returns combined stdout/stderr
```

## Usage Examples

### Session Management

```go
import "github.com/instruqt/git-exec/pkg/git"

// Create a new session with persistent configuration
// NewSession automatically detects existing repositories or creates new ones
session, err := git.NewSession("/path/to/project",
    git.WithUser("Jane Developer", "jane@instruqt.com"),
    git.WithInstruqtMetadata("user-123", "session-456", time.Now()),
    git.WithMetadata("track", "git-basics"),
    git.WithMetadata("environment", "production"),
    git.WithMetadata("project", "web-app"),
    git.WithMetadata("team", "frontend"),
)

// All operations automatically use the session's user context
err = session.Add([]string{"file.txt"})
err = session.Commit("Fix authentication bug") // Automatically uses Jane Developer as author

// Load an existing session (e.g., after service restart)
session, err = git.LoadSession("/path/to/project")
config := session.GetConfig()
fmt.Printf("Session ID: %s, User: %s\n", config.SessionID, config.UserName)
fmt.Printf("Project: %s, Environment: %s\n", 
    config.Metadata["project"], config.Metadata["environment"])

// Clone a repository using session context
// Create session first, then clone into it
cloneSession, err := git.NewSession("/path/to/clone-target",
    git.WithUser("John Dev", "john@instruqt.com"),
    git.WithInstruqtMetadata("user-789", "session-012", time.Now()),
    git.WithMetadata("source", "github-clone"),
)
// NewSession creates empty directory, remove .git so we can clone
os.RemoveAll("/path/to/clone-target/.git") 
err = cloneSession.Clone("https://github.com/user/repo", "/path/to/clone-target")

// Validate a session
err = git.ValidateSession("/path/to/project")
if err != nil {
    log.Printf("Invalid session: %v", err)
}

// Get session information
info, err := git.GetSessionInfo("/path/to/project")
fmt.Printf("User ID: %s, Session ID: %s\n", info.UserID, info.SessionID)

// Update user information in a session
err = session.UpdateUser("Jane Smith", "jane.smith@instruqt.com")

// Destroy session-specific configuration
err = session.Destroy()
```

### Basic Repository Operations

```go
// Create a new git instance (without session management)
git, err := git.NewGit()
if err != nil {
    log.Fatal(err)
}

// Initialize a repository
err = git.Init("/path/to/repo", InitWithBare())

// Clone a repository
err = git.Clone("https://github.com/user/repo", "/local/path",
    CloneWithBranch("main"),
    WithAuth(githubToken))

// Set working directory for subsequent operations
git.SetWorkingDirectory("/local/path")
```

### Working with Files

```go
// Add files (session automatically applies user context)
err = session.Add([]string{"*.go", "README.md"})

// Check status
files, err := session.Status(StatusWithPorcelain())
for _, file := range files {
    fmt.Printf("%s: %s\n", file.Path, file.Status)
}

// Commit changes (session user context applied automatically)
err = session.Commit("Update documentation", CommitWithSignoff())

// Or use temporary config for one-off operations
err = git.Commit("Quick fix",
    commands.WithConfig("user.name", "Temp Worker"),
    commands.WithConfig("user.email", "temp@example.com"),
    CommitWithSignoff())

// Multiple config values at once
configs := map[string]string{
    "user.name":     "Jane Developer",
    "user.email":    "jane@example.com",
    "commit.gpgsign": "false",
}
err = git.Commit("Multi-config commit", commands.WithConfigs(configs))
```

### Branch Management

```go
// List branches
branches, err := git.ListBranches(BranchWithRemotes())
for _, branch := range branches {
    if branch.IsCurrent {
        fmt.Printf("* %s\n", branch.Name)
    } else {
        fmt.Printf("  %s\n", branch.Name)
    }
}

// Create and checkout branch
err = git.CreateBranch("feature-xyz", BranchWithTrack("origin/main"))
err = git.Checkout(CheckoutWithBranch("feature-xyz"))

// Push with upstream
err = git.Push(PushWithRemote("origin"), 
    PushWithBranch("feature-xyz"),
    PushWithSetUpstream())
```

### Remote Operations

```go
// Add remote
err = git.AddRemote("upstream", "https://github.com/upstream/repo")

// Fetch from remote
remotes, err := git.Fetch(
    FetchWithRemote("upstream"),
    FetchWithPrune(),
    FetchWithTags())

// Pull changes
result, err := git.Pull(PullWithRebase())
if result != nil && !result.Success {
    fmt.Println("Conflicts:", result.Conflicts)
}
```

### Inspection Operations

```go
// View commit history
logs, err := git.Log(
    LogWithMaxCount(10),
    LogWithOneline(),
    LogWithGraph())

for _, log := range logs {
    fmt.Printf("%s %s - %s\n", 
        log.Hash[:7], 
        log.Author.Name, 
        log.Message)
}

// Show specific commit
commit, err := git.Show("HEAD~1", ShowWithStat())

// View differences
diffs, err := git.Diff(
    DiffWithNameOnly(),
    DiffWithCommit("HEAD~1", "HEAD"))
```

## Testing

The library uses testify for assertions and mockery for mocks:

```go
func TestAdd(t *testing.T) {
    g := &git{path: "git", wd: "/test/repo"}
    
    // Test basic add
    err := g.Add([]string{"file.txt"})
    require.NoError(t, err)
    
    // Test with options
    err = g.Add([]string{"*.go"}, 
        AddWithDryRun(),
        AddWithVerbose())
    require.NoError(t, err)
}
```

## Architecture

### Package Structure

```
pkg/git/
├── git.go                 # Main Git interface definition
├── command.go             # Command interface definition
├── commands/
│   ├── git.go            # Git implementation
│   ├── command.go        # Command execution logic
│   ├── parser.go         # Output parsing utilities
│   ├── add.go            # Add command implementation
│   ├── branch.go         # Branch operations
│   ├── checkout.go       # Checkout implementation
│   ├── clone.go          # Clone implementation
│   ├── commit.go         # Commit implementation
│   ├── diff.go           # Diff implementation
│   ├── fetch.go          # Fetch implementation
│   ├── init.go           # Init implementation
│   ├── log.go            # Log implementation
│   ├── merge.go          # Merge implementation
│   ├── misc.go           # Misc commands (config, reflog, etc.)
│   ├── pull.go           # Pull implementation
│   ├── push.go           # Push implementation
│   ├── rebase.go         # Rebase implementation
│   ├── remote.go         # Remote management
│   ├── reset.go          # Reset implementation
│   ├── show.go           # Show implementation
│   ├── status.go         # Status implementation
│   └── tag.go            # Tag operations
├── types/
│   └── types.go          # Structured types for Git objects
├── errors/
│   └── errors.go         # Error types and handling
├── mocks/
│   └── command.go        # Mock implementations for testing
└── test/
    └── *_test.go         # Integration tests
```

### Command Execution Flow

1. User calls interface method (e.g., `git.Add(files, options...)`)
2. Implementation creates command with `newCommand("add", args...)`
3. Options are applied via `ApplyOptions(opts...)`
4. Command is executed via `Execute()` or `ExecuteCombined()`
5. Output is parsed into structured types
6. Results or errors are returned to user

## Future Enhancements

### High Priority

1. **Complete Stub Implementations**:
   - `Config()` get/set operations  
   - `Checkout()` output parsing
   - Enhanced `Rebase()` with conflict handling

### Medium Priority

1. **Bare Repository Support**:
   - Specific bare repository operations
   - Reference management without working directory

2. **Enhanced Authentication**:
   - Credential helper integration
   - Token refresh mechanisms
   - SSH key support

3. **Performance Optimizations**:
   - Command batching
   - Output streaming for large operations
   - Parallel operations support

### Low Priority

1. **Additional Commands**:
   - Stash operations
   - Bisect support
   - Archive creation
   - Bundle operations

2. **Advanced Features**:
   - Progress callbacks
   - Custom merge strategies
   - Hook management

## VCS Service Integration

The library is designed to support the Instruqt VCS service requirements:

**Repository Synchronization**: Full support for clone, fetch, and push with authentication.

**Publishing Workflow**: Add, commit, and merge operations with structured error handling.

**Branch Management**: Complete branch and tag operations for multi-branch support.

**Remote Management**: Full remote configuration for GitHub integration.

**User Attribution**: User context support via options for proper commit attribution.

The library provides a solid foundation for building reliable version control workflows with type safety, structured errors, and comprehensive Git functionality.