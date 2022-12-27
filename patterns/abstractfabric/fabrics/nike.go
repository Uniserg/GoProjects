package fabrics

import (
	. "serguni/patterns/abstractfabric/abstractproducts"
	. "serguni/patterns/abstractfabric/products"
)

type Nike struct {
}

func (n *Nike) MakeShoe() IShoe {

	shoe := NikeShoe{}

	shoe.SetLogo("nike")
	shoe.SetSize(14)

	return &shoe
}

func (n *Nike) MakeShirt() IShirt {
	shirt := NikeShirt{}

	shirt.SetLogo("nike")
	shirt.SetSize(14)
	return &shirt
}
