package controller

import (
	"github.com/agill17/namespace-manager/pkg/controller/autokill"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, autokill.Add)
}
