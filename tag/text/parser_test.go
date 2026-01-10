package text

import (
	"fmt"
	"log/slog"
	"testing"

	"github.com/dubbikins/envy/v2/tag"
)


func TestParse(t *testing.T) {

	type Example struct {
		Field string `env:"FOO|BAZ;key1='value',key2='value'"`
	}
	var expected = []string{"FOO","BAZ"}
	ex := &Example{}
	node := tag.NewRootNode(ex)
	node = node.Descendant(node.Value().Elem(), nil)
	slog.Info("Node", "value", node.Value())
	var field = node.Value().Type().Field(0)
	node = node.Descendant(node.Value().Field(0), &field)
	// slog.Info("Node", "value", node.)
	
	var err = Parse("env", &node,LexEnvironmentVariableTag)
	if err != nil {
		t.Fatal(err)
	}
	for i, have := range node.TagValues() {
		if have != expected[i] {
			t.Fatalf("expected node.TagValues[%d] to be '%s', but was '%s'", i, expected[i], have)
		}
	}
	fmt.Println(node.Options())
	var keyname = "key1"
	var expected_key = "value"
	if have_key := node.Options()[keyname]; have_key != expected_key {
		t.Fatalf("expected node.Options[%s] = %q but was %q",keyname, expected_key, have_key)
	}
	keyname = "key2"
	if have_key := node.Options()[keyname]; have_key != expected_key {
		t.Fatalf("expected node.Options[%s] = %q but was %q",keyname, expected_key, have_key)
	}
}