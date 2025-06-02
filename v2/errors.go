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
func DoesNotMatchError(field_name, match_expression, actual string) error {
	return fmt.Errorf("matching error:\nStruct Field: [%s]\nTrying to match expression %v but value was '%s'", field_name, match_expression, actual)
}

var TAG_VALIDATION_ERROR = errors.New("invalid field definition; no tag value or default")
var INVALID_TAG_SYNTAX_ERROR = errors.New("invalid tag definition; tag syntax {tagvalue} (default={value}|options=[...comma seperated values]|required) ")
