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
	// translate the agent map into an array of agents
	GenerateAgentArrayFromMap() []T
}

type IServer[T agent.IAgent[T]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[T]
	// exposes server methods to agents for messaging, etc
	agent.IExposedServerFunctions[T]
	// gives access to number of iteration in simulator
	GetIterations() int
	// the set of functions defining how a 'game loop' should run
	// starts the agents' messaging session
	// RunMessagingSession()
	// begins simulator
	Start()
}



type RoundRunner interface {
	RunRound()
	RunTurn()
}
