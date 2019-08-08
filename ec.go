// el: error logging
package el

import (
	"errors"
	"fmt"
	"io"
	"runtime"
	"strings"
)

type node struct {
	previous error
	info     string
	file     string
	line     int
	fn       string
	function *runtime.Func
}

// New error chain node.
func NewNode(err error, fn, info string, level int) error {
	pc, file, line, _ := runtime.Caller(level)
	return &node{err, info, file, line, fn, runtime.FuncForPC(pc)}
}

// Wrap error with info and stacked error(s).
func Wrap(err error, info string) error {
	if err == nil {
		if len(info) == 0 {
			return nil
		}
		err = errors.New(info)
		info = ""
	}
	return NewNode(err, "", info, 2)
}

// Wrap error with fn, info and stacked error(s).
func WrapFn(err error, fn, info string) error {
	if err == nil {
		return nil
	}
	return NewNode(err, fn, info, 2)
}

// Take the string argument as a format string to generate the info for Wrap().
func Wrapf(err error, format string, a ...interface{}) error {
	if err == nil {
		return nil
	}
	return NewNode(err, "", fmt.Sprintf(format, a...), 2)
}

// NewError() as errors.NewError() with caller's location.
func NewError(info string) error {
	return NewNode(nil, "", info, 2)
}

// NewError() as errors.NewError() with caller's location.
func NewErrorFn(fn, info string) error {
	return NewNode(nil, fn, info, 2)
}

// Take the string argument as a format string to generate the info for NewError().
func NewErrorf(format string, a ...interface{}) error {
	return NewNode(nil, "", fmt.Sprintf(format, a...), 2)
}

// Take the string argument as a format string to generate the info for NewError().
func NewErrorfFn(fn, format string, a ...interface{}) error {
	return NewNode(nil, fn, fmt.Sprintf(format, a...), 2)
}

const prefixBlank = "\n  "

// dump to buffer
func (e *node) dump(out *strings.Builder) *strings.Builder {
	e.Dump(out)
	return out
}

func (e *node) print(out io.Writer) {
	fileName, packageName, functionName := e.GetName()

	_, _ = out.Write([]byte(fmt.Sprintf("[%v] %v:%v: %v()", packageName, fileName, e.line, functionName)))
	if len(e.info) > 0 {
		for _, s := range strings.Split(e.info, "\n") {
			_, _ = out.Write([]byte(prefixBlank))
			_, _ = out.Write([]byte(s))
		}
	}

}

// Dump to output stream
func (e *node) Dump(out io.Writer) {
	e.print(out)

	var err error
	for err = e.previous; err != nil; {
		switch n := err.(type) {
		case *node:
			_, _ = out.Write([]byte("\n"))
			n.print(out)
			err = n.previous
		default:
			_, _ = out.Write([]byte(prefixBlank))
			_, _ = out.Write([]byte(err.Error()))
			return
		}
	}
}

func (e *node) Error() string {
	return e.dump(&strings.Builder{}).String()
}

func (e *node) Unwrap() (err error) {
	return e.previous
}

func (e *node) Info() string {
	return e.info
}

func (e *node) Location() (file string, line int, function *runtime.Func) {
	return e.file, e.line, e.function
}

func (e *node) GetName() (fileName, packageName, functionName string) {
	functionName = e.function.Name()
	if n := strings.Index(functionName, "."); n > 0 {
		packageName, functionName = functionName[:n], functionName[n+1:]
	}
	if e.fn != "" {
		if n := strings.LastIndex(functionName, "."); n > 0 {
			functionName = functionName[:n+1] + e.fn
		}
	}

	fileName = e.file
	if n := strings.LastIndex(fileName, "/"); n > 0 {
		/*
		   // FIXME: show main's package location
		   if packageName == "main" {
		       const src = "/src/"
		       if m := strings.LastIndex(fileName[:n], src); m >= 0 {
		           packageName = fileName[m+len(src) : n]
		       }
		   }*/

		fileName = fileName[n+1:]
	}
	return
}

func (e *node) GetLine() (line int) {
	line = e.line
	return
}
