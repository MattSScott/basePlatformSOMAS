package infra

import (
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
}
