package infra

import (
	"github.com/google/uuid"
)

type IAgent[T any] interface {
	IMessagingProtocol
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
	// TODO
	NotifyAgentInactive()
}
