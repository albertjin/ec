package yaml

import "io"

type line interface {
    Output(prefix string, w io.Writer)
    HasErr() bool
}

type Line struct {
    Format func(s string) string
    Key    string
    Line   line
}

func (l Line) MakeKey() string {
    return MakeKey(l.Format, l.Key)
}

func NewLine(format func(string) string, key string, value line) Line {
    return Line{format, key, value}
}
