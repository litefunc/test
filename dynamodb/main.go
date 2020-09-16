package main

import "test/dynamodb/internal"

func main() {
	table := "CoffeeShop"

	internal.CreateItem(table)

	internal.Retrieve(table)
}
