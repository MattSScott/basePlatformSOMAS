package infra

import (
	"sync"

	"github.com/google/uuid"
)

type IAgent[T any] interface {
	// composes messaging passing capabilities
	iUnexportedAgentMethods
	IMessagingProtocol
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
	NotifyAgentInactive()
	//listenOnChannel(chan IMessage, chan ServerNotification, *sync.WaitGroup)
}

type iUnexportedAgentMethods interface {
	listenOnChannel(chan IMessage, chan ServerNotification, *sync.WaitGroup)
}
