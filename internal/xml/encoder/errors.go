package encoder

import "errors"

// ErrEmptyValue signals that a value is empty and should be serialized as an empty tag.
var ErrEmptyValue = errors.New("empty value")
