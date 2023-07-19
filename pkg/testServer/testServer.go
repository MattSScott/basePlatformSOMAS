package testserver

import (
	"fmt"
	baseServer "somas_base_platform/pkg/infra/server"
)

type MyServer struct {
	// embed struct for composition
	baseServer.BaseServer
	// keep second reference to superclass for unimplemented methods
	baseServer *baseServer.BaseServer
	name       string
	numFloors  int
}

func New() MyServer {
	return MyServer{
		baseServer: &baseServer.BaseServer{},
		numFloors:  20,
	}
}

func (ms *MyServer) Init() {
	ms.name = "My Server"
	fmt.Printf("Name field added as: %s \n", ms.name)
}
