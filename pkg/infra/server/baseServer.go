package infra

import (
	"fmt"
	"strconv"
	baseAgent "somas_base_platform/pkg/agents"
)


type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []baseAgent.Agent
}


func (bs *BaseServer) Init() {
	fmt.Println("Server Init")
	bs.NumAgents = 5
	bs.NumTurns = 4
	bs.Agents = make([]baseAgent.Agent, bs.NumAgents)
	for i := 0; i < bs.NumAgents; i++ {
		//converts the iteration to string
		name := strconv.Itoa(i)
		//creates new Agent
		bs.Agents[i]= baseAgent.NewAgent(name, 1, 100)
		
	}
}

func (bs *BaseServer) RunGameLoop(loopnum int) {
	fmt.Printf("Game Loop %d Running \n", loopnum)
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, element := range bs.Agents {
		fmt.Printf("agent %d has id: %s \n", index, element.ID)
		fmt.Printf("agent %d has name: %s \n", index, element.Name)
		fmt.Printf("agent %d has floor: %d \n", index, element.Floor)
		fmt.Printf("agent %d has energy: %d \n", index, element.Energy)
		fmt.Printf("_____________________________________________ \n")
		bs.Agents[index].UpdateAgent() 
		
		//TO DO: add the function for stages 


		
	}


}

func (bs *BaseServer) Start() {
	bs.Init()
	//LOOPS
	for i:=0; i < bs.NumTurns; i++{
		fmt.Printf("Loop: %d \n", i)
		bs.RunGameLoop(i)
	}
	
}
