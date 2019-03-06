package ecl_test

import (
    "testing"

    "github.com/albertjin/ecl"
)

func TestWrap(t *testing.T) {
    e0 := ecl.NewError("the original error\nbreak")
    e1 := ecl.Wrap(e0,"something wrong")
    e2 := ecl.Wrap(e1, "start")

    t.Log(e2.Error())

    e3 := ecl.Wrap(e2, "")
    t.Log(e3.Error())
}
