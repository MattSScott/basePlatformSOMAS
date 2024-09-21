package server

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/google/uuid"
)

type IAgentOperations[T agent.IAgent[T]] interface {
	// gives access to the agents in the simulator
	GetAgentMap() map[uuid.UUID]T
	// adds an agent to the server
	AddAgent(agentToAdd T)
	// removes an agent from the server
	RemoveAgent(agentToRemove T)
}

type IServer[T agent.IAgent[T]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[T]
	// exposes server methods to agents for messaging, etc
	agent.IExposedServerFunctions[T]
	// gives access to number of iteration in simulator
	GetIterations() int
	// begins simulator
	Start()
	//Signals the end of a messaging session. Either all agents send a message indicating they finished or server forcefully moves on after a set time period
	EndAgentListeningSession() bool
}

type GameRunner interface {
	RunStartOfIteration(int)
	RunTurn(int, int)
	RunEndOfIteration(int)
}
