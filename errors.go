package envy

import (
	"errors"
	"fmt"
)

var INVALID_FIELD_ERROR = errors.New("invalid field value")
var NOT_UNMARSHALLABLE_ERROR = errors.New("field is not unmarshallable")

func TagRequiredError(tagname, tag string) error {
	return fmt.Errorf("tag unmarshal error; (%s:%s) required but not set", tagname, tag)
}

func TagInvalidOptionError(tag, value string, options []string) error {
	return fmt.Errorf("tag error: '%s' options are %v but value is set to '%s'", tag, options, value)
}

var TAG_VALIDATION_ERROR = errors.New("invalid field definition; no tag value or default")
var INVALID_TAG_SYNTAX_ERROR = errors.New("invalid tag definition; tag syntax {tagvalue} (default={value}|options=[...comma seperated values]|required) ")
