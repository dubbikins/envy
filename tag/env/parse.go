package env
type Expander func(s string, mapping func(string) string) string
type Reader func(string) (string)



