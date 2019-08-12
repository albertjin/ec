package el

import (
	"testing"
)

func a() error {
	return Wrap(nil, "a")
}

func b() error {
	return Wrap(a(), "b")
}

type X struct {
}

func (x *X) Foo() error {
	return Wrap(a(), nil)
}

func (x X) Bar() error {
	return Wrap(a(), nil)
}

func TestWrap(t *testing.T) {
	t.Log("here")
	e0 := Wrap(nil, "the original error\nbreak")
	e1 := Wrap(e0, "something wrong")
	e2 := Wrap(e1, "start")

	t.Log("\n" + e2.Error())

	e3 := Wrap(e2, "")
	t.Log("\n" + e3.Error())

	t.Log("\n" + Wrap(ErrNil, "test").Error())
}

func TestWrapFn(t *testing.T) {
	t.Log("\n" + b().Error())
	x := &X{}
	t.Log("\n" + x.Foo().Error())
	t.Log("\n" + x.Bar().Error())
	t.Log("\n" + (func() error {
		return WrapFn(a(), "hello", "info")
	}()).Error())
}
