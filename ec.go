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
   Location() (file string, line int)
}

type node struct {
    info string
    previous error
    file string
    line int
}

// New error chain node.
func NewNode(info string, previous error, level int) Node {
    _, file, line, _ := runtime.Caller(level)
    return &node{info, previous, file, line}
}

// Wrap error with info and stacked error(s).
func Wrap(info string, err error) error {
    if err != nil {
        return NewNode(info, err, 2)
    }
    return nil
}

// Implement error.Error().
func (e *node) Error() string {
    return fmt.Sprintf("%v:%v: %v\n", e.file, e.line, e.info) + e.previous.Error()
}

func (e *node) Previous() error {
    return e.previous
}

func (e *node) Info() string {
    return e.info
}

func (e *node) Location() (file string, line int) {
    return e.file, e.line
}
