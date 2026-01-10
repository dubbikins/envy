package types_test

import (
	"os"
	"testing"

	"github.com/dubbikins/envy/v2"
	"github.com/dubbikins/envy/v2/types"
)

type ComplexConfig struct {
	Timeout types.Duration `env:"TIMEOUT" default:"5m" options:"1m,2m,5m,10m,15m"`
	String string `env:"STRING" default:"string"`
	Inner Config
}

type Config struct {
	Timeout types.Duration `env:"TIMEOUT" default:"5m" options:"1m,2m,5m,10m,15m"`
}
type StringConversion struct {
	String string `env:"STRING" default:"string"`
}

func BenchmarkComplexUnmarshal(b *testing.B) {
	// Set environment variable for benchmark
	os.Setenv("TIMEOUT", "2m")
	var cfg = &ComplexConfig{}
	for i := 0; i < b.N; i++ {
		envy.Unmarshal(cfg)
	}
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
	b.Run("simple unmarshal", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			envy.Unmarshal(strcv)
		}
	})
	
	b.ReportAllocs()
	if strcv.String != "test" {
		b.Fatalf("expected %s, but was %s", "test", strcv.String)
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
