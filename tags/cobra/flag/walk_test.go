package flag

import (
	"testing"

	"github.com/dubbikins/envy/v2/tag"
	"github.com/spf13/cobra"
)


func TestWalkFlags(t *testing.T) {
	cmd := &cobra.Command{

	}
	type Example struct {
		Name string `flag:"n|name;description='This is the description'"`
	}
	var have = &Example{}
	var err error
	if err = tag.Walk(tag.Chain(tag.UnmarshalText, TagWalkFn(cmd)), have); err != nil {
		t.Fatal(err)
	}
}