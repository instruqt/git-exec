package commands

// Add adds file contents to the index
func (g *git) Add(files []string, opts ...Option) error {
	args := []string{}
	
	if len(files) > 0 {
		args = append(args, files...)
	} else {
		args = append(args, ".")
	}
	
	cmd := g.newCommand("add")
	cmd.args = append(cmd.args, args...)
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
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