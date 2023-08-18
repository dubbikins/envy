package envy

import "time"

// Custom implementation for Parsing Durations from a string instead of parsing it as an int64
type Duration struct {
	time.Duration
}

func (d *Duration) Unmarshal(f *FieldReflection) (err error) {
	ref := f.Ref().(*Duration)
	tag, err := f.Tag()
	if err != nil {
		return err
	}
	ref.Duration, err = time.ParseDuration(tag.Value)
	if err != nil {
		return
	}
	f.Set(ref)
	return
}
