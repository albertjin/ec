package ec_test

import (
    "testing"

    "albertjin/ec"
)

func TestWrap(t *testing.T) {
    e0 := ec.NewError("the original error\nbreak")
    e1 := ec.Wrap(e0,"something wrong")
    e2 := ec.Wrap(e1, "start")

    t.Log(e2.Error())

    e3 := ec.Wrap(e2, "")
    t.Log(e3.Error())
}
