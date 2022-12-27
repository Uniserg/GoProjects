package main

import (
	"fmt"
	"serguni/patterns/abstractfabric"
	"serguni/patterns/abstractfabric/abstractproducts"
	"serguni/patterns/fabric"
)

func main() {

	fmt.Println("************* Using factory pattern *************")
	ak47, _ := fabric.GetGun("ak47")
	musket, _ := fabric.GetGun("musket")

	printDetails(ak47)
	printDetails(musket)

	fmt.Println("************* Using abstract factory pattern *************")

	adidasFactory, _ := abstractfabric.GetSportsFactory("adidas")
	nikeFactory, _ := abstractfabric.GetSportsFactory("nike")

	nikeShoe := nikeFactory.MakeShoe()
	nikeShirt := nikeFactory.MakeShirt()

	adidasShoe := adidasFactory.MakeShoe()
	adidasShirt := adidasFactory.MakeShirt()

	printShoeDetails(nikeShoe)
	printShirtDetails(nikeShirt)

	printShoeDetails(adidasShoe)
	printShirtDetails(adidasShirt)
}

func printDetails(g fabric.IGun) {
	fmt.Printf("Gun: %s", g.GetName())
	fmt.Println()
	fmt.Printf("Power: %d", g.GetPower())
	fmt.Println()
}

func printShoeDetails(s abstractproducts.IShoe) {
	fmt.Printf("Logo: %s", s.GetLogo())
	fmt.Println()
	fmt.Printf("Size: %d", s.GetSize())
	fmt.Println()
}

func printShirtDetails(s abstractproducts.IShirt) {
	fmt.Printf("Logo: %s", s.GetLogo())
	fmt.Println()
	fmt.Printf("Size: %d", s.GetSize())
	fmt.Println()
}
