package ec_test

import (
    "strings"
    "testing"

    "github.com/albertjin/ec"
)

func TestWrap(t *testing.T) {
    e0 := ec.NewError("the original error")
    e1 := ec.Wrap("something wrong", e0)

    // Note that e0 and e1 are called at line 11 and 12.
    s := e1.Error()
    if !strings.Contains(s, "github.com/albertjin/ec/ec_test.go:11: github.com/albertjin/ec_test.TestWrap()") ||
        !strings.Contains(s, "github.com/albertjin/ec/ec_test.go:12: github.com/albertjin/ec_test.TestWrap()") {
        t.Errorf("not expected:\n%v", s)
    }
}
