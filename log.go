package ec

import (
    "encoding/json"
    "fmt"
    "io"
    "regexp"
    "strings"
)

type logLine struct {
    key    string
    format func(s string) string
    line   LogLine
}

func (l logLine) makeKey() string {
    return logMakeKey(l.format, l.key)
}

type logText string

func (t logText) HasErr() bool {
    return false
}

func (t logText) Output(prefix string, w io.Writer) {
    if string(t) != "" {
        _, _ = w.Write([]byte(" "))
        logOutputText(prefix, string(t), w)
    }
}

type logError struct {
    err error
}

func (le logError) HasErr() bool {
    return true
}

func (le logError) Output(prefix string, w io.Writer) {
    var linePrefix []byte
    for e := le.err; e != nil; {
        var info string
        switch x := e.(type) {
        case Node:
            fileName, packageName, functionName := x.GetName()

            if linePrefix == nil {
                linePrefix = []byte("\n" + prefix + "- ")
                _, _ = w.Write([]byte("- "))
            } else {
                _, _ = w.Write(linePrefix)
            }
            _, _ = fmt.Fprintf(w, "[%v, %v, %v, %v()]", packageName, fileName, x.GetLine(), functionName)

            info, e = strings.TrimSpace(x.Info()), x.Previous()

        default:
            info, e = strings.TrimSpace(e.Error()), nil
        }

        if len(info) > 0 {
            if linePrefix == nil {
                linePrefix = []byte("\n" + prefix + "- ")
                _, _ = w.Write([]byte("- "))
            } else {
                _, _ = w.Write(linePrefix)
            }
            logOutputText(prefix, info, w)
        }
    }
}

type LogLine interface {
    Output(prefix string, w io.Writer)
    HasErr() bool
}

type LogKind int

const (
    LogMap  = LogKind(1)
    LogList = LogKind(2)
)

type Log struct {
    lines []logLine
    kind  LogKind

    errCount int
}

func NewLogMap() *Log {
    return &Log{kind: LogMap}
}

func NewLogList() *Log {
    return &Log{kind: LogList}
}

func (log *Log) KV(format func(string) string, key, value string) {
    log.lines = append(log.lines, logLine{key, format, logText(value)})
}

func (log *Log) V(value string) {
    log.lines = append(log.lines, logLine{"", nil, logText(value)})
}

func (log *Log) E(format func(string) string, key string, err error) {
    log.lines = append(log.lines, logLine{key, format, logError{err}})
    log.errCount++
}

func (log *Log) PushMap(format func(string) string, key string) *Log {
    n := &Log{kind: LogMap}
    log.lines = append(log.lines, logLine{key, format, n})
    return n
}

func (log *Log) PushList(format func(string) string, key string) *Log {
    n := &Log{kind: LogList}
    log.lines = append(log.lines, logLine{key, format, n})
    return n
}

func (log *Log) HasErr() bool {
    if log.errCount > 0 {
        return true
    }

    for _, line := range log.lines {
        if line.line.HasErr() {
            return true
        }
    }
    return false
}

func (log *Log) Output(prefix string, w io.Writer) {
    p0 := prefix
    p4 := p0 + "    "
    p2 := p4[:len(p0)+2]
    br := []byte("\n")

    switch log.kind {
    case LogMap:
        for i, line := range log.lines {
            if i > 0 {
                _, _ = w.Write([]byte(p0))
            }
            if key := line.makeKey(); len(key) > 0 {
                _, _ = w.Write([]byte(key))
            } else {
                _, _ = w.Write([]byte("'':"))
            }

            switch x := line.line.(type) {
            case *Log:
                var a, b string
                switch x.kind {
                case LogMap:
                    a, b = p2, " {}\n"
                case LogList:
                    a, b = p0, " []\n"
                }

                if len(x.lines) > 0 {
                    _, _ = w.Write(br)
                    _, _ = w.Write([]byte(a))
                    x.Output(a, w)
                } else {
                    _, _ = w.Write([]byte(b))
                }
            case logError:
                _, _ = w.Write(br)
                _, _ = w.Write([]byte(p0))
                x.Output(p0, w)
                _, _ = w.Write(br)
            default:
                x.Output(p0, w)
                _, _ = w.Write(br)
            }
        }
    case LogList:
        for i, line := range log.lines {
            key := line.makeKey()

            switch x := line.line.(type) {
            case *Log:
                var a, b, c string
                switch x.kind {
                case LogMap:
                    a, b, c = p4, " {}\n", "{}\n"
                case LogList:
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
                        x.Output(a, w)
                    } else {
                        _, _ = w.Write([]byte(b))
                    }
                } else {
                    if len(x.lines) > 0 {
                        x.Output(p2, w)
                    } else {
                        _, _ = w.Write([]byte(c))
                    }
                }
            case logError:
                if i > 0 {
                    _, _ = w.Write([]byte(p0))
                }
                _, _ = w.Write([]byte("- "))
                if len(key) > 0 {
                    _, _ = w.Write([]byte(key))
                    _, _ = w.Write(br)
                    _, _ = w.Write([]byte(p2))
                }
                x.Output(p2, w)
                _, _ = w.Write(br)
            default:
                if i > 0 {
                    _, _ = w.Write([]byte(p0))
                }
                if len(key) > 0 {
                    _, _ = w.Write([]byte("- "))
                    _, _ = w.Write([]byte(key))
                    x.Output(p2, w)
                } else {
                    _, _ = w.Write([]byte("-"))
                    x.Output(p0, w)
                }
                _, _ = w.Write(br)
            }
        }
    }
}

func logEscapeKey(s string) string {
    if logReText.MatchString(s) {
        e, _ := json.Marshal(s)
        return string(e)
    }

    return s
}

func logMakeKey(format func(string) string, key string) string {
    if len(key) == 0 {
        return key
    }

    key = logEscapeKey(key) + ":"
    if format == nil {
        return key
    }

    return format(key)
}

var logReText = regexp.MustCompile("(^[ \\t@|>\\-?!{}[\\]&`%*=,~\":'])|[#\\n]|:[ \\t]|([ \\t:]$)")

func logOutputText(prefix string, value string, w io.Writer) {
    if logReText.MatchString(value) {
        p := value
        if e := len(value) - 1; value[e] == '\n' {
            _, _ = w.Write([]byte("|2+\n"))
            p = value[:e]
        } else {
            _, _ = w.Write([]byte("|2-\n"))
        }

        for len(p) > 0 {
            _, _ = w.Write([]byte(prefix))
            _, _ = w.Write([]byte("  "))

            j := strings.IndexRune(p, '\n') + 1
            if j == 0 {
                _, _ = w.Write([]byte(p))
                break
            }
            _, _ = w.Write([]byte(p[:j]))
            p = p[j:]
        }
    } else {
        _, _ = w.Write([]byte(value))
    }
}
