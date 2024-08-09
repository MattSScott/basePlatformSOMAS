package infra

import (
	"github.com/google/uuid"
)

type IAgentOperations[T IAgent[T]] interface {
	// gives access to the agents in the simulator
	GetAgentMap() map[uuid.UUID]T
	// adds an agent to the server
	AddAgent(agentToAdd T)
	// removes an agent from the server
	RemoveAgent(agentToRemove T)
	// translate the agent map into an array of agents
	GenerateAgentArrayFromMap() []T
}

type IServer[T IAgent[T]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[T]
	// exposes server methods to agents for messaging, etc
	IExposedServerFunctions[T]
	// gives access to number of iteration in simulator
	GetIterations() int
	// the set of functions defining how a 'game loop' should run
	RunGameLoop()
	// starts the agents' messaging session
	RunMessagingSession()
	// begins simulator
	Start()
}

type IMessagingProtocol[T any] interface {
	SendSynchronousMessage(IMessage[T], []uuid.UUID)
	SendMessage(IMessage[T], []uuid.UUID)
	ReadChannel(uuid.UUID) <-chan IMessage[T]
	AcknowledgeClosure(uuid.UUID)
	AcknowledgeServerMessageReceived()
	// send notification that agent stopped talking session
	agentStoppedTalking(uuid.UUID)
}

type IExposedServerFunctions[T any] interface {
	IMessagingProtocol[T]
	// return hashset of all agent IDs
	ViewAgentIdSet() map[uuid.UUID]struct{}
	// return exposed functions for agent
	// AccessAgentByID(uuid.UUID) T
}

type RoundRunner interface {
	RunRound()
	RunTurn()
}

type ServerNotification int

const (
	StopListeningSpinner ServerNotification = iota
	StartListeningNotification
	EndListeningNotification
)
