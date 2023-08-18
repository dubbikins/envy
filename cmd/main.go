package main

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"

	"github.com/dubbikins/envy"
)

type MaxRuntime struct {
	time.Duration ``
}

func (d *MaxRuntime) UnmarshalText(value []byte) (err error) {
	d.Duration, err = time.ParseDuration(string(value))
	return
}

type MaxAmounts []string

func (p *MaxAmounts) UnmarshalText(text []byte) error {
	whole_str := string(text)
	*p = strings.Split(whole_str, ",")
	return nil
}

type Test struct {
	MaxRuntime MaxRuntime `env:"duration;default=5m"`
	MaxAmount  MaxAmounts `env:"amount;default=1,2,3"`
}
type Field struct {
	Value string `env:"test"`
}
type Message struct {
	Text string `json:"text"`
	Mode string `json:"mode"`
}
type Readme struct {
	Content string
}

func main() {
	test, err := envy.New(envy.FromEnvironmentAs[Test])
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(test)
	fmt.Println(string(data))
}
