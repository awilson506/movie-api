package api

import (
	"strconv"
)

// ErrorMessageContainer - hold some error messages
type ErrorMessageContainer struct {
	Errors map[string]string
}

// ValidateOptionalStringParam - confirm our params are ints
func ValidateOptionalStringParam(paramName string, param string, msg *ErrorMessageContainer) (string, bool) {

	if param == "" {
		return param, true
	}
	_, err := strconv.Atoi(param)

	if err != nil {
		msg.Errors[paramName] = "Please enter a valid value for: " + paramName
		return "", false
	}

	return param, true
}
