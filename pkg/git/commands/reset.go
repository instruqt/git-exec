package commands

// Reset resets current HEAD to the specified state
func (g *git) Reset(files []string, opts ...Option) error {
	cmd := g.newCommand("reset")
	
	if len(files) > 0 {
		cmd.args = append(cmd.args, files...)
	}
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Reset-specific options

// ResetWithSoft resets only HEAD
func ResetWithSoft() Option {
	return WithArgs("--soft")
}

// ResetWithMixed resets HEAD and index (default)
func ResetWithMixed() Option {
	return WithArgs("--mixed")
}

// ResetWithHard resets HEAD, index and working tree
func ResetWithHard() Option {
	return WithArgs("--hard")
}

// ResetWithMerge resets HEAD, index and working tree but keeps unmerged files
func ResetWithMerge() Option {
	return WithArgs("--merge")
}

// ResetWithKeep resets HEAD but keeps local changes
func ResetWithKeep() Option {
	return WithArgs("--keep")
}

// ResetWithRecurseSubmodules resets submodules too
func ResetWithRecurseSubmodules() Option {
	return WithArgs("--recurse-submodules")
}

// ResetWithNoRecurseSubmodules don't reset submodules
func ResetWithNoRecurseSubmodules() Option {
	return WithArgs("--no-recurse-submodules")
}

// ResetWithQuiet suppresses output
func ResetWithQuiet() Option {
	return WithArgs("--quiet")
}

// ResetWithCommit resets to a specific commit
func ResetWithCommit(commit string) Option {
	return WithArgs(commit)
}

// ResetWithPathspec resets specific pathspec
func ResetWithPathspec(pathspec string) Option {
	return WithArgs("--", pathspec)
}