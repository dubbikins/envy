package envy

type OptionsFunc[Options any] func(Options) error

func New[T any](options ...OptionsFunc[*T]) (*T, error) {
	var o *T = new(T)
	for _, option := range options {
		if err := option(o); err != nil {
			return o, err
		}
	}
	return o, nil
}

func FromEnvironment[T any](t *T) error {
	return Unmarshal(t)
}

func WithUnmarshalled[T any](do func(ptr *T)) {
	var t *T = new(T)
	if err := Unmarshal(t); err != nil {
		panic(err)
	}
	do(t)
}
