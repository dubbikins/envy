package tag

// An Error describes a failure to parse a struct tag
// and gives the offending expression.

// An SyntaxError describes a failure to parse a struct tag
type SyntaxError string

func (e SyntaxError) Error() string {
	return e.String()
}

const (
	// Unexpected error
	ErrInternalError SyntaxError = "envy/tag: internal error"
	ErrNestingDepth          SyntaxError = "expression nests too deeply"
	ErrLarge                 SyntaxError = "expression too large"
)

func (e SyntaxError) String() string {
	return string(e)
}
