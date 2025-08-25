# Git-Exec Examples

This directory contains comprehensive examples demonstrating the git-exec library functionality.

## Examples Overview

### 01_basic_operations.go
**Basic Git Operations**
- Repository initialization
- User configuration
- File creation and staging
- Status checking
- Committing changes
- Viewing commit history

**Key concepts**: Core Git workflow, command chaining, error handling

### 02_session_management.go
**Session Management with Persistent User Context**
- Creating sessions with user configuration
- Loading existing sessions
- Updating user information
- Multi-user collaboration
- Session validation and metadata
- Automatic user context application

**Key concepts**: Persistent sessions, metadata system, user management, collaboration

### 03_merge_operations.go
**Merge Operations and Conflict Resolution**
- Successful merges without conflicts
- Conflict detection and analysis
- Conflict resolution strategies (UseOurs, UseTheirs, Custom)
- Merge abort functionality
- Advanced merge options
- Comprehensive merge result information

**Key concepts**: Branch merging, conflict handling, resolution strategies, merge workflows

### 04_branch_management.go
**Branch Operations and Management**
- Listing and creating branches
- Switching between branches
- Making commits on different branches
- Branch-specific commit history
- Merging branches
- Cleaning up merged branches

**Key concepts**: Branch workflows, feature development, branch lifecycle management

## Running the Examples

Each example is self-contained and can be run independently:

```bash
# Run basic operations example
go run examples/01_basic_operations.go

# Run session management example  
go run examples/02_session_management.go

# Run merge operations example
go run examples/03_merge_operations.go

# Run branch management example
go run examples/04_branch_management.go
```

## Prerequisites

- Go 1.21.5 or later
- Git installed and available in PATH
- Write permissions for temporary directory creation

## Example Structure

Each example follows this pattern:

1. **Setup**: Creates temporary directories and initializes repositories
2. **Demonstration**: Shows specific functionality with real Git operations
3. **Explanation**: Provides context and explains what's happening
4. **Cleanup**: Automatically cleans up temporary files
5. **Summary**: Lists key concepts demonstrated

## Key Features Demonstrated

- **Type-safe Git operations** with structured error handling
- **Functional options pattern** for flexible command configuration
- **Comprehensive result types** with detailed operation information
- **Session-based user management** with persistent configuration
- **Advanced merge capabilities** with conflict resolution
- **Production-ready workflows** suitable for real applications

## Integration Examples

These examples show how git-exec can be used in:

- **CI/CD pipelines** with automated Git operations
- **Code review systems** with branch and merge management
- **Development tools** with session-based user context
- **Multi-tenant applications** with isolated Git operations
- **Educational platforms** with safe Git operation sandboxing

## Error Handling

All examples demonstrate proper error handling patterns:

- Checking errors from every Git operation
- Using structured error types
- Graceful failure and cleanup
- Informative error messages

## Testing

Examples create isolated temporary directories and are safe to run multiple times. They do not affect your local Git configuration or repositories.