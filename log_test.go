package ec

import (
    "bytes"
    "testing"
)

func TestEscapeKey(t *testing.T) {
    s := "\\'\""
    t.Log(s)
    t.Log(logEscapeKey(s))
}

func TestEscapeText(t *testing.T) {

    facts := map[string]string{
        ":":    "|2-\n  :",
        "a:":   "|2-\n  a:",
        "a:b":  "a:b",
        ":a:":  "|2-\n  :a:",
        "a: b": "|2-\n  a: b",

        "#":  "|2-\n  #",
        "a#": "|2-\n  a#",

        "|a": "|2-\n  |a",
        "a|": "a|",
    }

    for k, v := range facts {
        var b bytes.Buffer
        logOutputText("", k, &b)
        if b.String() != v {
            t.Error(b.String())
            t.Error(k)
            t.Error(v)
        }
    }
}
