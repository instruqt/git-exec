# Git Exec

Go library that wraps Git commands and returns structured output.

## Features

- Structured output via Go structs instead of string parsing
- Standard Go error handling with context
- Session management with persistent user configuration
- Merge conflict detection and resolution
- Mockery-compatible interfaces

## Installation

```bash
go get github.com/instruqt/git-exec
```

## Quick Start

```go
package main

import (
    "fmt"
    "log"
    "os"
    
    "github.com/instruqt/git-exec/pkg/git"
)

func main() {
    // Create Git instance
    gitInstance, err := git.NewGit()
    if err != nil {
        log.Fatal(err)
    }
    
    // Initialize repository
    err = gitInstance.Init("/path/to/repo")
    if err != nil {
        log.Fatal(err)
    }
    
    gitInstance.SetWorkingDirectory("/path/to/repo")
    
    // Configure user
    err = gitInstance.Config("user.name", "Your Name")
    if err != nil {
        log.Fatal(err)
    }
    
    err = gitInstance.Config("user.email", "you@example.com")
    if err != nil {
        log.Fatal(err)
    }
    
    // Add and commit files
    err = gitInstance.Add([]string{"."})
    if err != nil {
        log.Fatal(err)
    }
    
    err = gitInstance.Commit("Initial commit")
    if err != nil {
        log.Fatal(err)
    }
    
    // Check status
    files, err := gitInstance.Status()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Printf("Repository has %d files\n", len(files))
}
```

## Usage Examples

### Basic Operations

```go
// Get commit history
logs, err := gitInstance.Log(git.LogWithMaxCount("5"))
if err != nil {
    log.Fatal(err)
}

for _, logEntry := range logs {
    fmt.Printf("%s: %s - %s\n", logEntry.Hash[:8], logEntry.Message, logEntry.Author)
}

// Check repository status
files, err := gitInstance.Status()
if err != nil {
    log.Fatal(err)
}

for _, file := range files {
    fmt.Printf("%s %s\n", file.Status, file.Name)
}
```

### Session Management

Sessions maintain user configuration across operations:

```go
// Create session with user configuration
session, err := git.NewSession("/path/to/project",
    git.WithSessionUser("Alice Developer", "alice@company.com"),
    git.WithMetadata("project", "web-app"),
    git.WithMetadata("team", "frontend"),
)
if err != nil {
    log.Fatal(err)
}

// User context is automatically applied to all operations
err = session.Add([]string{"README.md"})
if err != nil {
    log.Fatal(err)
}

err = session.Commit("Initial project setup")
if err != nil {
    log.Fatal(err)
}

// Load existing session
loadedSession, err := git.LoadSession("/path/to/project")
if err != nil {
    log.Fatal(err)
}
```

### Branch Management

```go
// Create and switch to new branch
err = gitInstance.CreateBranch("feature/auth")
if err != nil {
    log.Fatal(err)
}

err = gitInstance.Checkout(git.CheckoutWithBranch("feature/auth"))
if err != nil {
    log.Fatal(err)
}

// List all branches
branches, err := gitInstance.ListBranches()
if err != nil {
    log.Fatal(err)
}

for _, branch := range branches {
    marker := " "
    if branch.Active {
        marker = "*"
    }
    fmt.Printf("  %s %s\n", marker, branch.Name)
}
```

### Merge Operations with Conflict Resolution

```go
// Attempt merge
result, err := gitInstance.Merge(git.MergeWithBranch("feature/auth"))
if err != nil {
    log.Fatal(err)
}

if result.Success {
    fmt.Printf("Merge successful: %s\n", result.MergedBranch)
} else {
    fmt.Printf("Merge conflicts detected in %d files\n", len(result.ConflictedFiles))
    
    // Show conflict details
    for _, conflict := range result.Conflicts {
        fmt.Printf("Conflict in %s:\n", conflict.Path)
        for i, section := range conflict.Sections {
            fmt.Printf("  Section %d: Our=%q, Their=%q\n", 
                i+1, section.OurContent, section.TheirContent)
        }
    }
    
    // Option 1: Resolve conflicts programmatically
    resolutions := []types.ConflictResolution{
        {
            FilePath: "conflicted-file.txt",
            UseOurs:  true, // or provide custom Resolution content
        },
    }
    
    err = gitInstance.ResolveConflicts(resolutions)
    if err != nil {
        log.Fatal(err)
    }
    
    // Continue merge
    err = gitInstance.MergeContinue()
    if err != nil {
        log.Fatal(err)
    }
}
```

#### Manual Conflict Resolution

Alternatively, conflicts can be resolved by manually editing files:

```go
// Attempt merge
result, err := gitInstance.Merge(git.MergeWithBranch("feature/auth"))
if err != nil {
    log.Fatal(err)
}

if !result.Success && len(result.Conflicts) > 0 {
    fmt.Printf("Merge conflicts detected in %d files:\n", len(result.ConflictedFiles))
    
    for _, conflictFile := range result.ConflictedFiles {
        fmt.Printf("  - %s\n", conflictFile)
    }
    
    fmt.Println("\nPlease resolve conflicts manually:")
    fmt.Println("1. Edit the conflicted files to resolve conflicts")
    fmt.Println("2. Remove conflict markers (<<<<<<, ======, >>>>>>)")
    fmt.Println("3. Stage resolved files and continue merge")
    
    // Wait for manual resolution...
    // User edits files externally in IDE/editor
    
    // After manual resolution, stage the resolved files
    err = gitInstance.Add(result.ConflictedFiles)
    if err != nil {
        log.Fatal(err)
    }
    
    // Continue the merge
    err = gitInstance.MergeContinue()
    if err != nil {
        log.Fatal(err)
    }
    
    fmt.Println("Merge completed successfully")
}
```

## Testing

Run tests:
```bash
go test ./...
```

Generate mocks for testing:
```bash
go install github.com/vektra/mockery/v2@latest
mockery --dir=pkg/git --name=Git
mockery --dir=pkg/git --name=Session
```

## Examples

See the `/examples` directory for complete working examples:

- `01_basic_operations.go` - Basic Git operations
- `02_session_management.go` - Session management with persistent context
- `03_merge_operations.go` - Merge operations and conflict resolution
- `04_branch_management.go` - Branch creation, switching, and management

Run examples:
```bash
go run examples/01_basic_operations.go
```