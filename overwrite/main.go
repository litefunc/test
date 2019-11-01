package main

import (
	"fmt"
	"test/overwrite/internal"
)

type Foo struct {
}

func (Foo) Call() {
	fmt.Println("Foo Called")
}

type Bar struct {
	Foo
}

type Baz struct {
	Foo
}

func (Baz) Call() {
	fmt.Println("Baz Called")
}

type B struct {
	internal.A
}

func (b B) Foo(internal.S) {
	fmt.Println("B.Foo()")
}

func main() {
	Foo{}.Call() // prints "Foo Called"
	Bar{}.Call() // prints "Foo Called"
	Baz{}.Call() // prints "Baz Called"

	b := B{A: internal.A{}}
	// b.Foo()
	internal.P(b)
}
