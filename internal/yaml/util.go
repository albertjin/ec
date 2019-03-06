package yaml

import (
    "encoding/json"
    "io"
    "regexp"
    "strings"
    "sync"
)

var reText *regexp.Regexp
var reOnce sync.Once

func GetReText() *regexp.Regexp {
    if reText != nil {
        return reText
    }

    reOnce.Do(func() {
        reText = regexp.MustCompile("(^[ \\t@|>\\-?!{}[\\]&`%*=,~\":'])|[#\\n]|:[ \\t]|([ \\t:]$)")
    })
    return reText
}

func EscapeKey(s string) string {
    if GetReText().MatchString(s) {
        e, _ := json.Marshal(s)
        return string(e)
    }

    return s
}

func MakeKey(format func(string) string, key string) string {
    if len(key) == 0 {
        return key
    }

    key = EscapeKey(key) + ":"
    if format == nil {
        return key
    }

    return format(key)
}

func OutputText(prefix string, value string, w io.Writer) {
    if GetReText().MatchString(value) {
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
