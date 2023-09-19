package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"

	"github.com/google/uuid"
)

type IAgent interface {
	message.IAgentMessenger[IAgent]
	Activity()
	UpdateAgent()
	GetID() uuid.UUID
}
