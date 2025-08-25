package git

// Add adds file contents to the index
func (g *gitImpl) Add(files []string, opts ...Option) error {
	args := []string{}
	
	if len(files) > 0 {
		args = append(args, files...)
	} else {
		args = append(args, ".")
	}
	
	cmd := g.newCommand("add")
	cmd.AddArgs(args...)
	
	// Apply all provided options
	cmd.ApplyOptions(opts...)
	
	_, err := cmd.Execute()
	return err
}