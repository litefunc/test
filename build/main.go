package main

type in interface {
	get()
}

func get(i in) {
	i.get()
}

func main() {
	var i int
	i.get()

}
