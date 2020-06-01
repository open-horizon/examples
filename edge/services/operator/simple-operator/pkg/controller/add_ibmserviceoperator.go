package controller

import (
	"projects/simple-operator/pkg/controller/ibmserviceoperator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, ibmserviceoperator.Add)
}
