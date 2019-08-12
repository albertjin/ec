package el

import "errors"

var (
	ErrNil  = errors.New("the pointer is unexpectedly nil")
	ErrVoid = errors.New("no descriptive information is available")
)
