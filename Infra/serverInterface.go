package infra

import "github.com/google/uuid"

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
	// gives access to number of iteration in simulator
	GetIterations() int
	// the set of functions defining how a 'game loop' should run
	RunGameLoop()
	// starts the agents' messaging session
	RunMessagingSession()
	// begins simulator
	Start()
}

type MessagingProtocolinterface interface {
	SendMessage(IMessage, uuid.UUID)
	ReadChannel(uuid.UUID) <-chan IMessage
	AcknowledgeClosure(uuid.UUID)
	AcknowledgeServerMessageReceived()
}

type IExposedServerFunctions interface {
	MessagingProtocol
	ViewAgentIdSet() map[uuid.UUID]struct{}
	getAgentServerChannel() *chan uuid.UUID
	agentStoppedTalking(uuid.UUID)
}

type RoundRunner interface {
	RunTurn()
}   

type ServerNotification int

const (
	StopListeningSpinner ServerNotification = iota
	StartListeningNotification
	EndListeningNotification
)
