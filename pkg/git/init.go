package git

// Init initializes a new Git repository
func (g *gitImpl) Init(path string, opts ...Option) error {
	cmd := g.newCommand("init", path)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}