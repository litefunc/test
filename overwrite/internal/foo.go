package internal

import "fmt"

type A struct {
}

func (a A) Foo(S) {
	fmt.Println("A.Foo()")
}

func (a A) Bar() {
	fmt.Println("A.Bar()")
}

type Client interface {
	Foo(S)
	Bar()
}

func P(f Client) {
	f.Foo(S{})
	f.Bar()
}

type S struct {
	A int
}
