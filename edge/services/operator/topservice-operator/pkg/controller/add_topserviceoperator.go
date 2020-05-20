package controller

import (
	"github.ibm.com/kube-operator/topservice-operator/pkg/controller/topserviceoperator"
)

func init() {
	// AddToManagerFuncs is a list of functions to create controllers and add them to a manager.
	AddToManagerFuncs = append(AddToManagerFuncs, topserviceoperator.Add)
}
