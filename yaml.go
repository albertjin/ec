package el

import (
    "io"

    "github.com/albertjin/el/internal/yaml"
)

type YamlKind int

const (
    YamlKindMap  YamlKind = 1
    YamlKindList YamlKind = 2
)

type YamlLog struct {
    lines []yaml.Line
    kind  YamlKind

    errCount int
}

func YamlLogMap() *YamlLog {
    return &YamlLog{kind: YamlKindMap}
}

func YamlLogList() *YamlLog {
    return &YamlLog{kind: YamlKindList}
}

func (l *YamlLog) KV(key interface{}, value string) {
    l.lines = append(l.lines, yaml.NewText(key, value))
}

func (l *YamlLog) V(value string) {
    l.lines = append(l.lines, yaml.NewText(nil, value))
}

func (l *YamlLog) E(key interface{}, err error) {
    l.lines = append(l.lines, yaml.NewError(key, err))
    l.errCount++
}

func (l *YamlLog) PushMap(key interface{}) *YamlLog {
    n := &YamlLog{kind: YamlKindMap}
    l.lines = append(l.lines, yaml.NewLine(key, n))
    return n
}

func (l *YamlLog) PushList(key interface{}) *YamlLog {
    n := &YamlLog{kind: YamlKindList}
    l.lines = append(l.lines, yaml.NewLine(key, n))
    return n
}

func (l *YamlLog) HasErr() bool {
    if l.errCount > 0 {
        return true
    }

    for _, line := range l.lines {
        if line.Line.HasErr() {
            return true
        }
    }
    return false
}

func (l *YamlLog) Output(context interface{}, prefix string, w io.Writer) {
    p0 := prefix
    p4 := p0 + "    "
    p2 := p4[:len(p0)+2]
    br := []byte("\n")

    switch l.kind {
    case YamlKindMap:
        for i, line := range l.lines {
            if i > 0 {
                _, _ = w.Write([]byte(p0))
            }
            if key := line.MakeKey(context); len(key) > 0 {
                _, _ = w.Write([]byte(key))
            } else {
                _, _ = w.Write([]byte("'':"))
            }

            switch x := line.Line.(type) {
            case *YamlLog:
                var a, b string
                switch x.kind {
                case YamlKindMap:
                    a, b = p2, " {}\n"
                case YamlKindList:
                    a, b = p0, " []\n"
                }

                if len(x.lines) > 0 {
                    _, _ = w.Write(br)
                    _, _ = w.Write([]byte(a))
                    x.Output(context, a, w)
                } else {
                    _, _ = w.Write([]byte(b))
                }
            case yaml.Error:
                _, _ = w.Write(br)
                _, _ = w.Write([]byte(p0))
                x.Output(context, p0, w)
                _, _ = w.Write(br)
            default:
                x.Output(context, p0, w)
                _, _ = w.Write(br)
            }
        }
    case YamlKindList:
        for i, line := range l.lines {
            key := line.MakeKey(context)

            switch x := line.Line.(type) {
            case *YamlLog:
                var a, b, c string
                switch x.kind {
                case YamlKindMap:
                    a, b, c = p4, " {}\n", "{}\n"
                case YamlKindList:
                    a, b, c = p2, " []\n", "[]\n"
                }

                if i > 0 {
                    _, _ = w.Write([]byte(p0))
                }
                _, _ = w.Write([]byte("- "))
                if len(key) > 0 {
                    _, _ = w.Write([]byte(key))
                    if len(x.lines) > 0 {
                        _, _ = w.Write(br)
                        _, _ = w.Write([]byte(a))
                        x.Output(context, a, w)
                    } else {
                        _, _ = w.Write([]byte(b))
                    }
                } else {
                    if len(x.lines) > 0 {
                        x.Output(context, p2, w)
                    } else {
                        _, _ = w.Write([]byte(c))
                    }
                }
            case yaml.Error:
                if i > 0 {
                    _, _ = w.Write([]byte(p0))
                }
                _, _ = w.Write([]byte("- "))
                if len(key) > 0 {
                    _, _ = w.Write([]byte(key))
                    _, _ = w.Write(br)
                    _, _ = w.Write([]byte(p2))
                }
                x.Output(context, p2, w)
                _, _ = w.Write(br)
            default:
                if i > 0 {
                    _, _ = w.Write([]byte(p0))
                }
                if len(key) > 0 {
                    _, _ = w.Write([]byte("- "))
                    _, _ = w.Write([]byte(key))
                    x.Output(context, p2, w)
                } else {
                    _, _ = w.Write([]byte("-"))
                    x.Output(context, p0, w)
                }
                _, _ = w.Write(br)
            }
        }
    }
}
