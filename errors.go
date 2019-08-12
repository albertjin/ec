package el

import "errors"

var (
	ErrNil        = errors.New("the pointer is nil")
	ErrUnexpected = errors.New("the value is unexpected")
	ErrVoid       = errors.New("no descriptive information is available")
)
