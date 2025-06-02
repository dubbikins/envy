package tag





type lexer struct {
	
}



type Example struct {
	Foo string `env:"FOO|FOOBAR" required:"true" default:"bar"`
	Bar string `env:"BAR" required:"true" default:"{{.Foo}}"`
}



