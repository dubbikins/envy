package envy

import "time"

// Custom implementation for Parsing Durations from a string instead of parsing it as an int64
type Duration struct {
	time.Duration
}

func (d *Duration) UnmarshalText(data []byte) (err error) {
	if len(data) == 0 {
		return nil
	}
	d.Duration, err = time.ParseDuration(string(data))
	return
}
