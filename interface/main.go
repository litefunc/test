package main

type a interface {
	a()
}

type b struct{}

func (rec b) a() {

}

func pa(a) {

}

func p([]a) {

}

func main() {
	var bs []b
	pa(b{})
	p(bs)
}
