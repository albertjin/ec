// error chain
package ec

import (
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

// implementation of error.Error()
func (e *node) Error() string {
    return e.info + " (" + e.previous.Error() + ")"
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
