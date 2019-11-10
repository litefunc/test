package pointer

import "testing"

type A struct {
	n int
}

func (r A) Set(n int) {
	r.n = n
}

type B struct {
	n int
}

func (r *B) Set(n int) {
	r.n = n
}

func TestValueReceiver(t *testing.T) {

	var a A
	a.Set(1)
	if want := 0; a.n != want {
		t.Errorf(`want:%d, got:%d`, want, a.n)
	}

	a1 := &A{}
	a1.Set(1)
	if want := 0; a1.n != want {
		t.Errorf(`want:%d, got:%d`, want, a1.n)
	}
}

func TestPointerReceiver(t *testing.T) {

	var b B
	b.Set(1)
	if want := 1; b.n != want {
		t.Errorf(`want:%d, got:%d`, want, b.n)
	}

	b1 := &B{}
	b1.Set(1)
	if want := 1; b1.n != want {
		t.Errorf(`want:%d, got:%d`, want, b1.n)
	}
}
