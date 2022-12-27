package fabrics

import (
	. "serguni/patterns/abstractfabric/abstractproducts"
	. "serguni/patterns/abstractfabric/products"
)

type Adidas struct {
}

func (a *Adidas) MakeShoe() IShoe {

	shoe := AdidasShoe{}

	shoe.SetLogo("adidas")
	shoe.SetSize(14)

	return &shoe
}

func (a *Adidas) MakeShirt() IShirt {
	shirt := AdidasShirt{}
	shirt.SetLogo("adidas")
	shirt.SetSize(14)

	return &shirt
}
