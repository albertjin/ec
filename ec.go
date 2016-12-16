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
    info string
    previous error
    file string
    line int
    function *runtime.Func
}

// New error chain node.
func NewNode(info string, previous error, level int) Node {
    pc, file, line, _ := runtime.Caller(level)
    return &node{info, previous, file, line, runtime.FuncForPC(pc)}
}

// Wrap error with info and stacked error(s).
func Wrap(info string, err error) error {
    if err == nil {
        return nil
    }
    return NewNode(info, err, 2)
}

// Take the string argument as a format string to generate the info for Wrap().
func Wrapf(format string, err error, a ...interface{}) error {
    if err == nil {
        return nil
    }
    return NewNode(fmt.Sprintf(format, a...), err, 2)
}

// NewError() as errors.NewError() with caller's location.
func NewError(info string) error {
    return NewNode(info, nil, 2)
}

// Take the string argument as a format string to generate the info for NewError().
func NewErrorf(format string, a ...interface{}) error {
    return NewNode(fmt.Sprintf(format, a...), nil, 2)
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
