package env

import (
	"log/slog"
	"os"
	"strings"

	"github.com/dubbikins/envy/v2/tag"
)
var environ_cache = map[string] map[string]string{}
/*
WalkFn returns a [TagWalkFunc] for handling the `env` struct tag

A FieldWriter is a function with the signature `func(node Node) error` and the behavior of WalkFn is as follows:

If node.Field.Tag does not contain the `env` struct tag, then return the value of calling next on node.
If node.Field.Tag contains a value v for the `env` struct tag, the value is handle as follows:
	Tag Syntax: `env_vars(;flags)`

	`env_vars` is treated as a list of environment variable names separated by the '|'
	Whitespace characters, as defined by Unicode, surrounding environment variable names are ignored.
	The first non-empty environment variable value encountered is written to node.

	If env_vars is "-" then return nil which indicats that further processing by next should not happen

	`flags` are specified by including a ';' followed by 0 or flags seperated by ','
		The available flags are:
			- `expand`
			If this flags is specified, then [os.Expand] is called on the value of the first non-empty
			environment variable specified in env_vars. The mapping function is a higher order function created from calling [ExpandFn] on node.ctx which handles values specified with a '.' prefix will by looking up the value in the node's context
			otherwise values are looked up using os.Getenv

	ex: `env:"FOO|BAR"`
		when FOO=bar
		Then "bar" will be written to the node

	ex: `env:"FOO;expand`
		When
			BAZ=foo
			FOO=${BAZ}bar
		Then "foobar" will be written to the node

	ex: `env:"FOO;expand`
		When
			BAZ=bar
			FOO=${.BAZ}bar
		And node.Context.Value(".BAZ") returns "foo"
		Then "foobar" will be written to the node
*/
func WalkFn(next tag.WalkFn) tag.WalkFn {
	return func(node *tag.Node) (err error) {
		//If this isn't a struct field element, then we should stop walking the tag
		if node.Field() == nil { //&& node.Value().Kind() != reflect.Pointer
			return 
		}
		var value string
		if err = node.Parse("env", tag.LexEnvironmentVariableTag); err != nil || node.Skipped(){
			return 
		}
		var Getenv Reader = os.Getenv
		var expand = node.FlagIsSet("expand")
		if env_source, found := node.Options()[".env"]; found {
			var osFallback = strings.HasPrefix(env_source, "os+")
			if osFallback {
				env_source = strings.TrimPrefix(env_source, "os+")
			}
			env_source = strings.Replace(env_source, "/", string(os.PathSeparator), -1)
			if strings.HasPrefix(env_source, ".") {
				var cwd string
				if cwd, err = os.Getwd(); err != nil {
					return
				}
				env_source = strings.Replace(env_source, ".", cwd, 1)
			}
			var environ map[string]string
			var data []byte
			var found bool
			
			if environ, found = environ_cache[env_source]; !found {
				environ = map[string]string{}
				if data, err = os.ReadFile(env_source); err != nil {
					return 
				}
				var vars = strings.Split(string(data), "\n")
				for _, kv := range vars {
						if split := strings.IndexRune(string(kv), '='); split >= 0 {
							environ[string(kv[:split])] = string(kv[split+1:])
						}
				}
				environ_cache[env_source] = environ
				slog.Info("Using custom env reader", "source", "env_reader_source", "values", environ )
			}
			slog.Info("Using custom environ reader")
			Getenv = func(s string) (string) {
				var _found bool
				if s, _found = environ[s]; !_found && osFallback {
					return os.Getenv(s)
				}
				return s
			}
		}
		
		
		for _, envVar := range node.TagValues() {
			if value = Getenv(envVar); value != "" {
				slog.Info("Fetching Env Var", envVar, value)

				if expand {
					// Getenv = func(s string) string {
					// 	slog.Info("expanding env", "name", s)
					// 	return 
					// }
					value = os.Expand(value, ExpandFn(node, Getenv)) //ExpandFn(node, Getenv)
				}
				if  err = node.UnmarshalText([]byte(value)); err != nil {
					return
				}
				break
			}
		}
		node.Reset()
		return next(node)
	}
}