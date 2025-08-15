package commands

// Commit creates a new commit with the given message
func (g *git) Commit(message string, opts ...Option) error {
	cmd := g.newCommand("commit", "-m", message)
	
	// Apply quiet by default
	cmd.args = append(cmd.args, "-q")
	
	// Apply all provided options
	cmd.applyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
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