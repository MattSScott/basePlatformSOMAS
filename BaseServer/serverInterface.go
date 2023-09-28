package baseserver

import (
	baseagent "github.com/MattSScott/basePlatformSOMAS/BaseAgent"
	"github.com/google/uuid"
)

type IAgentOperations[T baseagent.IAgent[T]] interface {
	// gives access to the agents in the simulator
	GetAgentMap() map[uuid.UUID]T
	// adds an agent to the server
	AddAgent(agentToAdd T)
	// removes an agent from the server
	RemoveAgent(agentToRemove T)
	// translate the agent map into an array of agents
	GenerateAgentArrayFromMap() []T
}

type IServer[T baseagent.IAgent[T]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[T]
	// gives access to number of iteration in simulator
	GetIterations() int
	// the set of functions defining how a 'game loop' should run
	RunGameLoop()
	// begins simulator
	Start()
}
