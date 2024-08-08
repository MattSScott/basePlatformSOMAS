package infra

import (
	"sync"

	"github.com/google/uuid"
)

type IAgent[T any] interface {
	// composes necessary server functions for agent access
	IExposedServerFunctions[T]
	// returns the unique ID of an agent
	GetID() uuid.UUID
	// allows agent to update their internal state
	UpdateAgentInternalState()
	// TODO
	NotifyAgentInactive()
	// TODO: move to better location
	RunSynchronousMessaging()
	// allow agent to listen on channel
	listenOnChannel(chan IMessage, chan ServerNotification, *sync.WaitGroup)
}
