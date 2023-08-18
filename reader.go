package envy

import (
	"os"
	"sync"
)

var reader Reader
var lock = sync.Mutex{}

type Reader interface {
	Get(string) string
}

type OSEnvironmentReader struct {
}

func (g *OSEnvironmentReader) Get(key string) string {
	return os.Getenv(key)
}

func init() {
	lock.Lock()
	reader = &OSEnvironmentReader{}
	lock.Unlock()
}

func UseReader(r Reader) {
	lock.Lock()
	reader = r
	lock.Unlock()
}
