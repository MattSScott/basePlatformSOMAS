package baseagent

import (
	"github.com/google/uuid"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
)

type Agent interface {
	message.Agent
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
