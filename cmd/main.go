package main

import (
	"fmt"
	"os"
	"time"

	"github.com/dubbikins/envy"
)

type Duration struct {
	time.Duration
}

func (d *Duration) Unmarshal(f *envy.FieldReflection) (err error) {
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

type Test struct {
	Duration Duration `env:"duration;options=[5m,10m,15m]"`
}
type Field struct {
	Value string `env:"test"`
}

func main() {
	os.Setenv("test", "5s")
	os.Setenv("struct_field", "value")
	os.Setenv("duration", "5m")
	test := &Test{}

	err := envy.Unmarshal(test, &envy.OSEnvironmentReader{})
	if err != nil {
		panic(err)
	}
	fmt.Println(test)
}
