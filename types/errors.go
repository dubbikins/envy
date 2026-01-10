package types


type TypeConversionError string

func (e TypeConversionError) Error() string {
	return string(e)
}

type ContextError string

func (e ContextError) Error() string {
	return string(e)
}