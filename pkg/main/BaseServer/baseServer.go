package baseserver

import (
	"fmt"

	baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
	"github.com/google/uuid"
)

type BaseServer[T baseagent.IAgent[T]] struct {
	numTurns int
	agents   map[uuid.UUID]T
}

func (bs *BaseServer[T]) GetAgents() map[uuid.UUID]T {
	return bs.agents
}

func (bs *BaseServer[T]) AddAgent(agentToAdd T) {
	bs.agents[agentToAdd.GetID()] = agentToAdd
}

func (bs *BaseServer[T]) RemoveAgent(agentToAdd T) {
	delete(bs.agents, agentToAdd.GetID())
}

func (bs *BaseServer[T]) GetNumTurns() int {
	return bs.numTurns
}

func (bs *BaseServer[T]) RunGameLoop() {
	for _, agent := range bs.agents {
		fmt.Printf("Agent %s updating state \n", agent.GetID())
		agent.UpdateAgentInternalState()
	}
}

func (bs *BaseServer[T]) Start() {
	fmt.Printf("Server initialised with %d agents \n", len(bs.agents))
	fmt.Print("\n")
	//LOOPS
	for i := 0; i < bs.numTurns; i++ {
		fmt.Printf("Game Loop %d running... \n \n", i)
		fmt.Printf("Main game loop running... \n \n")
		bs.RunGameLoop()
		fmt.Printf("\nMain game loop finished. \n \n")
		fmt.Printf("Messaging session started... \n \n")
		bs.RunMessagingSession()
		fmt.Printf("\nMessaging session completed \n \n")
		fmt.Printf("Game Loop %d completed. \n", i)
	}

}

type AgentGenerator[T baseagent.IAgent[T]] func() T

type AgentGeneratorCountPair[T baseagent.IAgent[T]] struct {
	generator AgentGenerator[T]
	count     int
}

func MakeAgentGeneratorCountPair[T baseagent.IAgent[T]](generatorFunction AgentGenerator[T], count int) AgentGeneratorCountPair[T] {
	return AgentGeneratorCountPair[T]{
		generator: generatorFunction,
		count:     count,
	}
}

func (bs *BaseServer[T]) GenerateAgentArrayFromMap() []T {

	agentMapToArray := make([]T, len(bs.agents))

	i := 0
	for _, ag := range bs.agents {
		agentMapToArray[i] = ag
		i++
	}
	return agentMapToArray
}

func (bs *BaseServer[T]) RunMessagingSession() {

	agentArray := bs.GenerateAgentArrayFromMap()

	for _, agent := range bs.agents {
		allMessages := agent.GetAllMessages(agentArray)
		for _, msg := range allMessages {
			recipients := msg.GetRecipients()
			for _, recip := range recipients {
				if agent.GetID() == recip.GetID() {
					continue
				}
				msg.InvokeMessageHandler(recip)
			}
		}
	}
}

func (bs *BaseServer[T]) initialiseAgents(m []AgentGeneratorCountPair[T]) {

	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agent := pair.generator()
			bs.AddAgent(agent)
		}
	}

}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[T baseagent.IAgent[T]](mapper []AgentGeneratorCountPair[T], iterations int) *BaseServer[T] {
	serv := &BaseServer[T]{
		agents:   make(map[uuid.UUID]T),
		numTurns: iterations,
	}
	serv.initialiseAgents(mapper)
	return serv
}
