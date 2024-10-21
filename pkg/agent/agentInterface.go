package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/v2/internal/diagnosticsEngine"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/google/uuid"
)

type IExposedServerFunctions[T any] interface {
	// return hashset of all agent IDs
	ViewAgentIdSet() map[uuid.UUID]struct{}
	// return exposed functions for agent
	AccessAgentByID(uuid.UUID) T
	// allows base agent to deliver message
	DeliverMessage(message.IMessage[T], uuid.UUID)
	// notify that agent has completed talking phase
	AgentStoppedTalking(uuid.UUID)
	// return max number of threads spawnable by an agent
	GetAgentMessagingBandwidth() int
	// return diagnostic engine used for tracking message data
	GetDiagnosticEngine() diagnosticsEngine.IDiagnosticsEngine
}

type IMessagingFunctions[T any] interface {
	// allows for creation of a base message
	CreateBaseMessage() message.BaseMessage
	// allows for sending a message to a single recipient
	SendMessage(message.IMessage[T], uuid.UUID)
	// allows for sending a message to a single recipient synchronously
	SendSynchronousMessage(message.IMessage[T], uuid.UUID)
	// allows for sending an async message across the entire system
	BroadcastMessage(message.IMessage[T])
	// allows for sending a sync message across the entire system
	BroadcastSynchronousMessage(message.IMessage[T])
	// signals end of agent's listening session
	SignalMessagingComplete()
}

type IAgent[T any] interface {
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// exposes messaging functions for base agent
	IMessagingFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
}
