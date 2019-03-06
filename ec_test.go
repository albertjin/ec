package el_test

import (
    "testing"

    "github.com/albertjin/el"
)

func TestWrap(t *testing.T) {
    e0 := el.NewError("the original error\nbreak")
    e1 := el.Wrap(e0,"something wrong")
    e2 := el.Wrap(e1, "start")

    t.Log(e2.Error())

    e3 := el.Wrap(e2, "")
    t.Log(e3.Error())
}
