package envy

import (
	"os"
	"testing"
)

type Config struct {
	Timeout Duration `env:"TIMEOUT" default:"5m" options:"1m,2m,5m,10m,15m"`
}
type StringConversion struct {
	String string `env:"STRING" default:"string"`
}

func BenchmarkWithUnmarshalled(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("STRING", "test")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		WithUnmarshalled[Config](func(ptr *Config) {})
	}
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2m")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(FromEnvironment[Config])
	}
}
func BenchmarkStringUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("STRING", "test")
	strcv := &StringConversion{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(strcv)
	}
}

func BenchmarkCustomDurationNew(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2h15m")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		New(FromEnvironment[Config])
	}
}

func BenchmarkCustomDurationUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2h15m")
	cfg := &Config{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		Unmarshal(cfg)
	}
}

// func BenchmarkRemoveListPrefix(b *testing.B) {
// 	// Set environment variable for benchmark
// 	bytes := []byte("[test]")

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		removeListPrefix(bytes)
// 	}
// }

// func BenchmarkRemoveListPrefixWithTrim(b *testing.B) {
// 	// Set environment variable for benchmark
// 	str := "[test]"

// 	b.ResetTimer()
// 	for i := 0; i < b.N; i++ {
// 		strings.Split(strings.Trim(str, "[({})]"), ",")
// 	}

// }
