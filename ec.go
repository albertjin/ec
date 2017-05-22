// error chain
package ec

import (
    "fmt"
    "runtime"
)

// error chain node
type Node interface {
   error
   Previous() error
   Info() string
   Location() (file string, line int, function *runtime.Func)
}

type node struct {
    previous error
    info string
    file string
    line int
    function *runtime.Func
}

// New error chain node.
func NewNode(err error, info string, level int) Node {
    pc, file, line, _ := runtime.Caller(level)
    return &node{err, info, file, line, runtime.FuncForPC(pc)}
}

// Wrap error with info and stacked error(s).
func Wrap(err error, info string) error {
    if err == nil {
        return nil
    }
    return NewNode(err, info, 2)
}

// Take the string argument as a format string to generate the info for Wrap().
func Wrapf(err error, format string, a ...interface{}) error {
    if err == nil {
        return nil
    }
    return NewNode(err, fmt.Sprintf(format, a...), 2)
}

// NewError() as errors.NewError() with caller's location.
func NewError(info string) error {
    return NewNode(nil, info, 2)
}

// Take the string argument as a format string to generate the info for NewError().
func NewErrorf(format string, a ...interface{}) error {
    return NewNode(nil, fmt.Sprintf(format, a...), 2)
}

// Implement error.Error().
func (e *node) Error() string {
    if e.previous != nil {
        return fmt.Sprintf("[ec] %v:%v: %v() %v\n", e.file, e.line, e.function.Name(), e.info) + e.previous.Error()
    }
    return fmt.Sprintf("[ec] %v:%v: %v() %v", e.file, e.line, e.function.Name(), e.info)
}

func (e *node) Previous() error {
    return e.previous
}

func (e *node) Info() string {
    return e.info
}

func (e *node) Location() (file string, line int, function *runtime.Func) {
    return e.file, e.line, e.function
}
