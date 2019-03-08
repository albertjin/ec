package yaml

import "io"

type Text string

func (t Text) HasErr() bool {
    return false
}

func (t Text) Output(context interface{}, prefix string, w io.Writer) {
    if string(t) != "" {
        _, _ = w.Write([]byte(" "))
        OutputText(prefix, string(t), w)
    }
}

func NewText(key interface{}, value string) Line {
    return Line{key, Text(value)}
}
