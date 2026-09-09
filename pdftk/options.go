package pdftk

// Option adjusts an invocation before it runs.
type Option func(cmd *command)

// OptionExecutable invokes name instead of the default "pdftk" — "pdftk-java",
// for instance, or an absolute path.
func OptionExecutable(name string) Option {
	return func(cmd *command) {
		cmd.executable = name
	}
}

// OptionFlatten flattens the PDF before output.
func OptionFlatten() Option {
	return func(cmd *command) {
		cmd.args = append(cmd.args, "flatten")
	}
}
