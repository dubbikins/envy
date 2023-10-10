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

func FromEnvironmentAs[T any](t *T) error {
	return Unmarshal(t)
}
func FromEnvironmentWithOptionAs[T any](options ...func(tag TagMiddleware)) func(t *T) error {
	return func(t *T) error {
		return Unmarshal(t, options...)
	}
}
