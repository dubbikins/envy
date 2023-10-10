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
	MaxRuntime MaxRuntime `env:"duration" default:"5m"`
	MaxAmount  MaxAmounts `env:"amount" default:"1,2,3"`
	Message
	Field
	FieldPtr   *Field
	unexported string
}
type Field struct {
	Value string `env:"test;required"`
}
type Message struct {
	Text string `json:"text"`
	Mode string `json:"mode"`
}
type Readme struct {
	Content string
}

type Brad struct {
	Name string `env:"NAME" required:"true"`
}

func main() {
	//test := &Brad{}
	rm, err := envy.New(envy.FromEnvironmentAs[Brad])
	// err := envy.Unmarshal(test, func(mw envy.TagMiddleware) {
	// 	mw.Push(
	// 		func(next envy.TagHandler) envy.TagHandler {
	// 			return envy.TagHandlerFunc(func(ctx context.Context, field reflect.StructField) error {
	// 				fmt.Println("I'm in the custom tag parser")

	// 				return next.UnmarshalField(ctx, field)
	// 			})
	// 		})
	// })
	if err != nil {
		panic(err)
	}
	data, _ := json.Marshal(rm)
	fmt.Println(string(data))
}
