package server

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/google/uuid"
)

type IAgentOperations[T agent.IAgent[T]] interface {
	// gives access to the agents in the simulator
	GetAgentMap() map[uuid.UUID]T
	// adds an agent to the server
	AddAgent(T)
	// removes an agent from the server
	RemoveAgent(T)
}

type IGameStateController interface {
	// gives access to number of iteration in simulator
	GetIterations() int
	// gives access to number of turns in simulator
	GetTurns() int
	// injects a GameRunner interface into the server
	SetGameRunner(GameRunner)
	// begins simulator
	Start()
}

type GameRunner interface {
	RunStartOfIteration(int)
	RunTurn(int, int)
	RunEndOfIteration(int)
}

type IServer[T agent.IAgent[T]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[T]
	// exposes server methods to agents for messaging, etc
	agent.IExposedServerFunctions[T]
	// exposes server methods for controlling game state
	IGameStateController
	// toggle logging of messaging diagnostics to console (default false)
	ReportMessagingDiagnostics()
}
