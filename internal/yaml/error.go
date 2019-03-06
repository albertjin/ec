package yaml

import (
    "fmt"
    "io"
    "strings"
)

type Error struct {
    Err error
}

func (le Error) HasErr() bool {
    return true
}

type node interface {
    Previous() error
    Info() string
    GetName() (fileName, packageName, functionName string)
    GetLine() (line int)
}

func (le Error) Output(prefix string, w io.Writer) {
    var linePrefix []byte
    for e := le.Err; e != nil; {
        var info string
        if x, ok := e.(node); ok {
            fileName, packageName, functionName := x.GetName()

            if linePrefix == nil {
                linePrefix = []byte("\n" + prefix + "- ")
                _, _ = w.Write([]byte("- "))
            } else {
                _, _ = w.Write(linePrefix)
            }
            _, _ = fmt.Fprintf(w, "[%v, %v, %v, %v()]", packageName, fileName, x.GetLine(), functionName)

            info, e = strings.TrimSpace(x.Info()), x.Previous()
        } else {
            info, e = strings.TrimSpace(e.Error()), nil
        }

        if len(info) > 0 {
            if linePrefix == nil {
                linePrefix = []byte("\n" + prefix + "- ")
                _, _ = w.Write([]byte("- "))
            } else {
                _, _ = w.Write(linePrefix)
            }
            OutputText(prefix, info, w)
        }
    }
}

func NewError(format func(string) string, key string, value error) Line {
    return Line{format, key, Error{value}}
}
