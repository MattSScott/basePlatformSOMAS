package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"

	"github.com/google/uuid"
)

type IAgentNetworking[T any] interface {
	// returns the full agent network map
	GetNetwork() map[uuid.UUID]T
	// adds an agent object to the network
	AddAgentToNetwork(agent T)
	// removes an agent object from the network
	RemoveAgentFromNetwork(agent T)
}

type IAgent[T any] interface {
	// composes messaging passing capabilities
	message.IAgentMessenger[T]
	// handles network operations
	IAgentNetworking[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
}
