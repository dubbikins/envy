package envy

// var reader KeyValuer
// var lock = sync.Mutex{}

// type KeyValuer interface {
// 	Get(string) string
// }

// type OSEnvironmentReader struct {
// }

// func (g *OSEnvironmentReader) Get(key string) string {
// 	return os.Getenv(key)
// }

// func init() {
// 	lock.Lock()
// 	reader = &OSEnvironmentReader{}
// 	lock.Unlock()
// }

// func UseReader(r KeyValuer) {
// 	lock.Lock()
// 	reader = r
// 	lock.Unlock()
// }
