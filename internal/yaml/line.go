package yaml

import "io"

type line interface {
    Output(context interface{}, prefix string, w io.Writer)
    HasErr() bool
}

type Line struct {
    Key  interface{}
    Line line
}

func (l Line) MakeKey(context interface{}) string {
    return MakeKey(context, l.Key)
}

func NewLine(key interface{}, value line) Line {
    return Line{key, value}
}
