package script

import (
	"fmt"
)

type Label struct {
	Address int
	Name string
	Comment string
	FarLabel bool
}

func AutoLabel(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("L%04X", address),
	}
}

func AutoLabelVar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("Var_%04X", address),
	}
}

func AutoLabelFar(address int) *Label {
	return &Label{
		Address: address,
		Name: fmt.Sprintf("F%04X", address),
		FarLabel: true,
	}
}

func NewLabel(address int, name string) *Label {
	return &Label{
		Address: address,
		Name: name,
	}
}

