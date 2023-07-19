package testserver

import (
	"fmt"
	basePlatform "somas_base_platform/pkg/infra/server"
)

type MyServer struct {
	basePlatform *basePlatform.BaseServer
	name         string
	numFloors    int
}

func New() MyServer {
	return MyServer{
		basePlatform: &basePlatform.BaseServer{},
		name:         "My Server",
		numFloors:    20,
	}
}

func (ms *MyServer) Init() {
	// ms.name = "My Server"
	// ms.numFloors = 50
	ms.basePlatform.Init()
	fmt.Printf("Name field added as: %s \n", ms.name)
	fmt.Println(ms.basePlatform.Agents)
}

// func (ms *MyServer) Start() {
// 	ms.Init()
// 	ms.RunGameLoop()
// }

// func RunGameLoop() {

// }
