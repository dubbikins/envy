package envy

func Unmarshal(s any) (err error) {
	var r *Reflection
	r, err = NewReflection(s)
	if err != nil {
		return
	}
	for _, field := range r.Element.Fields {
		err = field.Unmarshal()
		if err != nil {
			return
		}
	}
	return err
}
func FromEnvironmentAs[T any](t *T) error { return Unmarshal(t) }
