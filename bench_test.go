package envy_test

import (
	"os"
	"testing"

	"github.com/dubbikins/envy"
)

type Config struct {
	Timeout envy.Duration `env:"TIMEOUT" default:"5m" options:"1m,2m,5m,10m,15m"`
}
type StringConversion struct {
	String string `env:"STRING" default:"string"`
}

func BenchmarkSimpleUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2m")
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envy.New(envy.FromEnvironment[Config])
	}
}
func BenchmarkStringUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("STRING", "test")
	strcv := &StringConversion{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envy.Unmarshal(strcv)
	}
}

func BenchmarkCustomDurationNew(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2h15m")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envy.New(envy.FromEnvironment[Config])
	}
}

func BenchmarkCustomDurationUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2h15m")
	cfg := &Config{}
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		envy.Unmarshal(cfg)
	}
}
