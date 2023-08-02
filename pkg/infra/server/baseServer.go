package infra

import (
	"strconv"
	"fmt"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
	agent1 "somas_base_platform/pkg/agents/AgentTypes/agent1"
 	agent2 "somas_base_platform/pkg/agents/AgentTypes/agent2"
	
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
//type AgentGenerator func() baseAgent.BaseAgent //* causes different error
//type AgentGenerator func() *agent2.Agent2
//type AgentGenerator func() *baseAgent.Agent
type AgentGenerator func() baseAgent.Agent

type AgentGeneratorCountPair struct {
	generator AgentGenerator
	count     int
}

func initAgents() []baseAgent.BaseAgent {
	m := make([]AgentGeneratorCountPair, 6)
	m[0] = AgentGeneratorCountPair{baseAgent.GetAgent, 4}
	m[1] = AgentGeneratorCountPair{agent2.GetAgent, 3}
	m[2] = AgentGeneratorCountPair{agent1.GetAgent, 2}

	agents := make([]baseAgent.BaseAgent, getNumAgents(m))
	agentCount := 0

	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agents[agentCount] = pair.generator()
			agentCount++
		}
	}

	return agents
}

func getNumAgents(pairs []AgentGeneratorCountPair) int {

	numAgents := 0

	for _, pair := range pairs {
		numAgents += pair.count
	}

	return numAgents
}