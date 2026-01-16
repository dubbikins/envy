package tag

import "fmt"

// An Error describes a failure to parse a struct tag
// and gives the offending expression.

type _Error string
// An SyntaxError describes a failure to parse a struct tag
type SyntaxError struct {
	_Error
}

func (e _Error) Error() string {
	return e.String()
}
func (e _Error) String() string {
	return string(e)
}


func ErrRequiredTagIsZero(node *Node) error {
	return fmt.Errorf("required field is zero: %s", node.Field().Name)
}


var (
	// Unexpected error
	ErrInternalError  = SyntaxError{"envy/tag: internal error"}
	ErrNestingDepth         =  SyntaxError {"expression nests too deeply"}
	ErrLarge                 =  SyntaxError {"expression too large"}
)


