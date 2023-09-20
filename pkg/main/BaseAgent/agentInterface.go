package baseagent

import (
	message "github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"

	"github.com/google/uuid"
)

type IAgent[T any] interface {
	message.IAgentMessenger[T]
	Activity()
	UpdateAgent()
	GetID() uuid.UUID
}
