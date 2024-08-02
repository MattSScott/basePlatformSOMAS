package infra

import "github.com/google/uuid"

type IAgentOperations[U IExposedAgentFunctions, T IAgent[U]] interface {
	// gives access to the agents in the simulator
	GetAgentMap() map[uuid.UUID]T
	// adds an agent to the server
	AddAgent(agentToAdd T)
	// removes an agent from the server
	RemoveAgent(agentToRemove T)
	// translate the agent map into an array of agents
	GenerateAgentArrayFromMap() []T
	// casts agent...
	CastAgentToExposedAgentFunctions(T) U
}

type IServer[U IExposedAgentFunctions, T IAgent[U]] interface {
	// gives operations for adding/removing agents from the simulator
	IAgentOperations[U, T]
	// gives access to number of iteration in simulator
	GetIterations() int
	// the set of functions defining how a 'game loop' should run
	RunGameLoop()
	// starts the agents' messaging session
	RunMessagingSession()
	// begins simulator
	Start()
}

type IMessagingProtocol interface {
	RunSynchronousMessaging()
	SendSynchronousMessage(IMessage, []uuid.UUID)
	SendMessage(IMessage, []uuid.UUID)
	ReadChannel(uuid.UUID) <-chan IMessage
	AcknowledgeClosure(uuid.UUID)
	AcknowledgeServerMessageReceived()
}

type IExposedServerFunctions[T IExposedAgentFunctions] interface {
	// return hashset of all agent IDs
	ViewAgentIdSet() map[uuid.UUID]struct{}
	// return exposed functions for agent
	GetAgentFromID(uuid.UUID) T
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
