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

// Wrap error with info and stacked error(s). When current previous error is nil,
// this method creates an error node with location.
func Wrap(info string, err error) error {
    return NewNode(info, err, 2)
}

func NewError(info string) error {
    return NewNode(info, nil, 2)
}

// Implement error.Error().
func (e *node) Error() string {
    if e.previous != nil {
        return fmt.Sprintf("%v:%v: %v\n", e.file, e.line, e.info) + e.previous.Error()
    }
    return fmt.Sprintf("%v:%v: %v", e.file, e.line, e.info)
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
