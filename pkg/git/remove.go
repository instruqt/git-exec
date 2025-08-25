package git

// Remove removes files from the working tree and from the index
func (g *gitImpl) Remove(opts ...Option) error {
	cmd := g.newCommand("rm")
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}