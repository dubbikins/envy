package envy

import (
	"errors"
	"fmt"
	"strings"
)

var INVALID_FIELD_ERROR = errors.New("invalid field value")
var NOT_UNMARSHALLABLE_ERROR = errors.New("field is not unmarshallable")

func RequiredError(tagname, tag_details string) error {
	return fmt.Errorf("tag unmarshal error; (%s) required but not set.\n %s", tagname, tag_details)
}
func InvalidOptionError(value string, options []string) error {
	return fmt.Errorf("tag unmarshalling error: options are [%s] but the value was set to '%s'", strings.Join(options, ","), value)
}
func DoesNotMatchError(tag, value string, matcher string) error {
	return fmt.Errorf("tag error: '%s' matcher expression is %v but is set to '%s'", tag, matcher, value)
}

var TAG_VALIDATION_ERROR = errors.New("invalid field definition; no tag value or default")
var INVALID_TAG_SYNTAX_ERROR = errors.New("invalid tag definition; tag syntax {tagvalue} (default={value}|options=[...comma seperated values]|required) ")
