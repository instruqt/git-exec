package git

// Commit creates a new commit with the given message
func (g *gitImpl) Commit(message string, opts ...Option) error {
	cmd := g.newCommand("commit", "-m", message)
	
	// Apply quiet by default
	cmd.AddArgs("-q")
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}