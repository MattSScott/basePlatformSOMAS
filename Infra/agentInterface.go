package infra

import (

	"github.com/google/uuid"
)

type IAgent[T any] interface {
	// composes messaging passing capabilities
	IAgentMessenger[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
}

