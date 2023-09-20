package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"

	"github.com/google/uuid"
)

type IAgent[T any] interface {
	// composes messaging passing capabilities
	message.IAgentMessenger[T]
	// allows agent to update their internal state
	UpdateAgentInternalState()
	// returns the unique ID of an agent
	GetID() uuid.UUID
}
