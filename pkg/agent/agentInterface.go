package agent

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type IAgent[T any] interface {
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
	// signals end of agent's listening session
	NotifyAgentFinishedMessaging()
	// allows for synchronous messaging to be run
	RunSynchronousMessaging()
	// allows for creation of a base message
	CreateBaseMessage() message.BaseMessage
}

type IMessagingProtocol[T any] interface {
	SendSynchronousMessage(message.IMessage[T], uuid.UUID)
	SendMessage(message.IMessage[T], uuid.UUID)
	BroadcastMessage(message.IMessage[T])
	AgentStoppedTalking(uuid.UUID)
}

type IExposedServerFunctions[T any] interface {
	IMessagingProtocol[T]
	// return hashset of all agent IDs
	ViewAgentIdSet() map[uuid.UUID]struct{}
	// return exposed functions for agent
	AccessAgentByID(uuid.UUID) T
}
