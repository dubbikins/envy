package envy

type Options struct {
	Middleware []Middleware
}

func (opts *Options) Unmarshal(s any) {

}

func DefaultMiddleware() []Middleware {
	return []Middleware{
		WithRequiredTag,
		WithMatchesTag,
		WithOptionsTag,
		WithEnvTag,
		WithDefaultTag,
	}
}

var default_options *Options

func init() {
	default_options = &Options{
		Middleware: DefaultMiddleware(),
	}
}
