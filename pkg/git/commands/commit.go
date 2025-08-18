package commands

import gitpkg "github.com/instruqt/git-exec/pkg/git"

// Commit creates a new commit with the given message
func (g *git) Commit(message string, opts ...gitpkg.Option) error {
	cmd := g.newCommand("commit", "-m", message)
	
	// Apply quiet by default
	cmd.AddArgs("-q")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}

// Commit-specific options

// CommitWithAuthor sets the author for the commit (alias for WithUser)
func CommitWithAuthor(name, email string) gitpkg.Option {
	return WithUser(name, email)
}

// CommitWithAll automatically stages all modified and deleted files
func CommitWithAll() gitpkg.Option {
	return WithArgs("--all")
}

// CommitWithAmend replaces the tip of the current branch
func CommitWithAmend() gitpkg.Option {
	return WithArgs("--amend")
}

// CommitWithNoEdit uses the previous commit message without launching an editor
func CommitWithNoEdit() gitpkg.Option {
	return WithArgs("--no-edit")
}

// CommitWithAllowEmpty allows creating a commit with no changes
func CommitWithAllowEmpty() gitpkg.Option {
	return WithArgs("--allow-empty")
}

// CommitWithAllowEmptyMessage allows a commit with an empty message
func CommitWithAllowEmptyMessage() gitpkg.Option {
	return WithArgs("--allow-empty-message")
}

// CommitWithSignoff adds a Signed-off-by line
func CommitWithSignoff() gitpkg.Option {
	return WithArgs("--signoff")
}

// CommitWithGPGSign signs the commit with GPG
func CommitWithGPGSign(keyid string) gitpkg.Option {
	if keyid == "" {
		return WithArgs("--gpg-sign")
	}
	return WithArgs("--gpg-sign=" + keyid)
}

// CommitWithNoVerify bypasses pre-commit and commit-msg hooks
func CommitWithNoVerify() gitpkg.Option {
	return WithArgs("--no-verify")
}