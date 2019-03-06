package yaml

import "io"

type Text string

func (t Text) HasErr() bool {
    return false
}

func (t Text) Output(prefix string, w io.Writer) {
    if string(t) != "" {
        _, _ = w.Write([]byte(" "))
        OutputText(prefix, string(t), w)
    }
}

func NewText(format func(string) string, key, value string) Line {
    return Line{format, key, Text(value)}
}
