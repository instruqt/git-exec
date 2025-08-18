package commands
import gitpkg "github.com/instruqt/git-exec/pkg/git"

// Reset resets current HEAD to the specified state
func (g *git) Reset(files []string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("reset")
	
	if len(files) > 0 {
		cmd.AddArgs(files...)
	}
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Reset-specific options

// ResetWithSoft resets only HEAD
func ResetWithSoft() gitpkg.Option {
	return WithArgs("--soft")
}

// ResetWithMixed resets HEAD and index (default)
func ResetWithMixed() gitpkg.Option {
	return WithArgs("--mixed")
}

// ResetWithHard resets HEAD, index and working tree
func ResetWithHard() gitpkg.Option {
	return WithArgs("--hard")
}

// ResetWithMerge resets HEAD, index and working tree but keeps unmerged files
func ResetWithMerge() gitpkg.Option {
	return WithArgs("--merge")
}

// ResetWithKeep resets HEAD but keeps local changes
func ResetWithKeep() gitpkg.Option {
	return WithArgs("--keep")
}

// ResetWithRecurseSubmodules resets submodules too
func ResetWithRecurseSubmodules() gitpkg.Option {
	return WithArgs("--recurse-submodules")
}

// ResetWithNoRecurseSubmodules don't reset submodules
func ResetWithNoRecurseSubmodules() gitpkg.Option {
	return WithArgs("--no-recurse-submodules")
}

// ResetWithQuiet suppresses output
func ResetWithQuiet() gitpkg.Option {
	return WithArgs("--quiet")
}

// ResetWithCommit resets to a specific commit
func ResetWithCommit(commit string) gitpkg.Option {
	return WithArgs(commit)
}

// ResetWithPathspec resets specific pathspec
func ResetWithPathspec(pathspec string) gitpkg.Option {
	return WithArgs("--", pathspec)
}