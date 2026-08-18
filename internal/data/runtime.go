package data

import (
	"errors"
	"fmt"
	"strconv"
	"strings"
)

var ErrorInvalidRuntimeFormat = errors.New("invalid runtime format")

type Runtime int

func (r Runtime) MarshalJSON() ([]byte, error) {
	json := fmt.Sprintf("%d mins", r)

	quotedJSON := strconv.Quote(json)

	return []byte(quotedJSON), nil
}

// UnmarshalJSON() modifies the receiver, so a pointer reference must be used.
// Otherwise, we will only be modifying a copy that gets discarded at the end of
// the method.
func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {
	unquotedJSONValue, err := strconv.Unquote(string(jsonValue))
	if err != nil {
		return ErrorInvalidRuntimeFormat
	}

	// Split the string to parse the number.
	parts := strings.Split(unquotedJSONValue, " ")
	if len(parts) != 2 || parts[1] != "mins" {
		return ErrorInvalidRuntimeFormat
	}

	runtime, err := strconv.Atoi(parts[0])
	if err != nil {
		return ErrorInvalidRuntimeFormat
	}

	*r = Runtime(runtime)

	return nil
}
