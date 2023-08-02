package infra

import (
	"fmt"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
	"strconv"
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []baseAgent.Agent
}

type AgentGenerator func() baseAgent.Agent

type AgentGeneratorCountPair struct {
	generator AgentGenerator
	count     int
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
		bs.Agents[i] = baseAgent.NewAgent(name)

	}
}

func (bs *BaseServer) RunGameLoop(loopnum int) {
	fmt.Printf("Game Loop %d Running \n", loopnum)
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, agent := range bs.Agents {
		fmt.Printf("agent %d  \n", index)
		fmt.Printf("_____________________________________________ \n")
		agent.UpdateAgent()

		//TO DO: add the function for stages

	}

}

func (bs *BaseServer) Start() {
	bs.Init()
	//LOOPS
	for i := 0; i < bs.NumTurns; i++ {
		fmt.Printf("Loop: %d \n", i)
		bs.RunGameLoop(i)
	}

}
