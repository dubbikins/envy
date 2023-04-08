package envy

import "os"

type EnvironmentReader interface {
	Getenv(string) string
}

type OSEnvironmentReader struct {
}

func (g *OSEnvironmentReader) Getenv(key string) string {
	return os.Getenv(key)
}
