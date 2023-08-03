package infra

import (
	"fmt"

	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []baseAgent.Agent
}

func (bs *BaseServer) RunGameLoop() {
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, agent := range bs.Agents {
		fmt.Printf("agent %d  \n", index)
		fmt.Printf("_____________________________________________ \n")
		agent.UpdateAgent()
		agent.Activity()
		//TO DO: add the function for stages

	}

}

func (bs *BaseServer) Init() {
	fmt.Println("Server initialised")
}

func (bs *BaseServer) Start() {
	//LOOPS
	for i := 0; i < bs.NumTurns; i++ {
		fmt.Printf("Game Loop %d Running \n", i)
		bs.RunGameLoop()
	}

}

type AgentGenerator func() baseAgent.Agent

type AgentGeneratorCountPair struct {
	generator AgentGenerator
	count     int
}

func MakeAgentGeneratorCountPair(generatorFunction AgentGenerator, count int) AgentGeneratorCountPair {
	return AgentGeneratorCountPair{
		generator: generatorFunction,
		count:     count,
	}
}

func (bs *BaseServer) initialiseAgents(m []AgentGeneratorCountPair) {

	agents := make([]baseAgent.Agent, getNumAgents(m))
	agentCount := 0

	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agents[agentCount] = pair.generator()
			agentCount++
		}
	}

	bs.Agents = agents
	bs.NumAgents = len(agents)
}

func getNumAgents(pairs []AgentGeneratorCountPair) int {

	numAgents := 0

	for _, pair := range pairs {
		numAgents += pair.count
	}

	return numAgents
}

func CreateServer(mapper []AgentGeneratorCountPair) *BaseServer {
	// generate the server and return it
	serv := &BaseServer{}
	serv.initialiseAgents(mapper)
	return serv
}
