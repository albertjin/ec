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
	info     interface{}
	file     string
	line     int
	fn       string
	function *runtime.Func
}

// New error chain node.
func newNode(err error, fn string, info interface{}, level int) error {
	if err == nil {
		switch v := info.(type) {
		case nil:
			err = ErrVoid
		case string:
			if len(v) == 0 {
				err = ErrVoid
			} else {
				err = errors.New(v)
			}
		case error:
			err = v
		default:
			err = errors.New(fmt.Sprintf("%v", info))
		}
		info = nil
	}
	pc, file, line, _ := runtime.Caller(level)
	return &node{err, info, file, line, fn, runtime.FuncForPC(pc)}
}

// Wrap error with info and stacked error(s).
func Wrap(err error, info interface{}) error {
	return newNode(err, "", info, 2)
}

// Wrap error with fn, info and stacked error(s).
func WrapFn(err error, fn string, info interface{}) error {
	return newNode(err, fn, info, 2)
}

// Take the string argument as a format string to generate the info for Wrap().
func Wrapf(err error, format string, a ...interface{}) error {
	return newNode(err, "", fmt.Sprintf(format, a...), 2)
}

// Take the string argument as a format string to generate the info for WrapFn().
func WrapFnf(fn, format string, a ...interface{}) error {
	return newNode(nil, fn, fmt.Sprintf(format, a...), 2)
}

const prefixBlank = "\n  "

// dump to buffer
func (e *node) dump(out *strings.Builder) *strings.Builder {
	e.Dump(out)
	return out
}

func (e *node) print(out io.StringWriter) {
	fileName, functionName := e.GetName()

	_, _ = out.WriteString(fmt.Sprintf("%v:%v: %v()", fileName, e.line, functionName))
	if e.info != nil {
		info, ok := e.info.(string)
		if !ok {
			info = fmt.Sprintf("%v", e.info)
		}
		if len(info) > 0 {
			for _, s := range strings.Split(info, "\n") {
				_, _ = out.WriteString(prefixBlank)
				_, _ = out.WriteString(s)
			}
		}
	}

}

// Dump to output stream
func (e *node) Dump(out io.StringWriter) {
	e.print(out)

	var err error
	for err = e.previous; err != nil; {
		switch n := err.(type) {
		case *node:
			_, _ = out.WriteString("\n")
			n.print(out)
			err = n.previous
		default:
			for _, s := range strings.Split(err.Error(), "\n") {
				_, _ = out.WriteString(prefixBlank)
				_, _ = out.WriteString(s)
			}
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

func (e *node) Info() interface{} {
	return e.info
}

func (e *node) Location() (file string, line int, function *runtime.Func) {
	return e.file, e.line, e.function
}

func (e *node) GetName() (fileName, functionName string) {
	functionName = e.function.Name()

	if e.fn != "" {
		if n := strings.LastIndex(functionName, "."); n > 0 {
			functionName = functionName[:n+1] + e.fn
		}
	}

	fileName = e.file
	return
}

func (e *node) GetLine() (line int) {
	line = e.line
	return
}
