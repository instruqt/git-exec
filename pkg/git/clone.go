package git

// Clone clones a repository
func (g *gitImpl) Clone(url, destination string, opts ...Option) error {
	cmd := g.newCommand("clone", url, destination)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}