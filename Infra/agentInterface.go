package infra

import (
	"github.com/google/uuid"
)

type IAgent[T IExposedAgentFunctions] interface {
	// composes messaging passing capabilities
	IMessagingProtocol
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
}

type IExposedAgentFunctions interface {
	IAgent[IExposedAgentFunctions]
}
