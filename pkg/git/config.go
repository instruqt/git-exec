package git

// Config sets a git configuration value
func (g *gitImpl) Config(key string, value string, opts ...Option) error {
	cmd := g.newCommand("config", key, value)
	cmd.ApplyOptions(opts...)
	_, err := cmd.Execute()
	return err
}