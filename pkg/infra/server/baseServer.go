package infra

import (
	"fmt"
	"strconv"
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []string
}

func (bs *BaseServer) Init() {
	fmt.Println("Server Init")
	bs.NumAgents = 5
	bs.NumTurns = 50
	bs.Agents = make([]string, bs.NumAgents)
	for i := 0; i < bs.NumAgents; i++ {
		bs.Agents[i] = strconv.Itoa(i)
	}
}

func (bs *BaseServer) RunGameLoop() {
	fmt.Println("Game Loop Running")
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, element := range bs.Agents {
		fmt.Printf("agent %d has id: %s \n", index, element)
	}

}

func (bs *BaseServer) Start() {
	bs.Init()
	bs.RunGameLoop()
}
