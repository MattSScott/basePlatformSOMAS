package infra

import (
	"fmt"
	"strconv"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
	Agent1 	  "somas_base_platform/pkg/agents/AgentTypes/agent1"
	Agent2    "somas_base_platform/pkg/agents/AgentTypes/agent2"
	
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []*baseAgent.BaseAgent
}

type AgentGenerator func() baseAgent.BaseAgentInterface

type AgentGeneratorCountPair struct {
	generator AgentGenerator
	count     int
}


func (bs *BaseServer) Init() {
	fmt.Println("Server Init")
	bs.NumAgents = 5
	bs.NumTurns = 4
	bs.Agents = make([]*baseAgent.BaseAgent, bs.NumAgents)
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
		//fmt.Printf("agent %d has id: %s \n", index, agent.ID)
		//fmt.Printf("agent %d has name: %s \n", index, agent.Name)
		//fmt.Printf("agent %d has floor: %d \n", index, element.Floor)
		//fmt.Printf("agent %d has energy: %d \n", index, element.Energy)
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
