package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
	"github.com/google/uuid"
)

type IAgent interface {
	message.IAgentMessaging
	Activity()
	UpdateAgent()
	GetID() uuid.UUID
	// Messaging
	// GetMsg() string
	// GetNet() []BaseAgent
	// GetRcv() []BaseAgent
	// SetMsg(s string)
	// SetNet(a []BaseAgent)
	// SetRcv(a []BaseAgent)
}
