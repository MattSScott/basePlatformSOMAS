package infra

import (
	"fmt"
	//agent1 "somas_base_platform/pkg/agents/AgentTypes/agent1"
	//agent2 "somas_base_platform/pkg/agents/AgentTypes/agent2"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []baseAgent.Agent
}

func (bs *BaseServer) Init(m []AgentGeneratorCountPair) {
	fmt.Println("Server Init")
	initialisedAgents := bs.initAgents(m)
	bs.Agents = initialisedAgents
	bs.NumAgents = len(initialisedAgents)
}

func (bs *BaseServer) RunGameLoop(loopnum int) {
	fmt.Printf("Game Loop %d Running \n", loopnum)
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, agent := range bs.Agents {
		fmt.Printf("agent %d  \n", index)
		fmt.Printf("_____________________________________________ \n")
		agent.UpdateAgent()
		agent.Activity()
		//TO DO: add the function for stages

	}

}

func (bs *BaseServer) Start(m []AgentGeneratorCountPair) {
	bs.Init(m)
	//LOOPS
	for i := 0; i < bs.NumTurns; i++ {
		fmt.Printf("Loop: %d \n", i)
		bs.RunGameLoop(i)
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

func (bs *BaseServer) initAgents(m []AgentGeneratorCountPair) []baseAgent.Agent {


	agents := make([]baseAgent.Agent, getNumAgents(m))
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
