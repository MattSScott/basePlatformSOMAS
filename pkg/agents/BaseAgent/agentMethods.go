package baseagent
import (
	
	"github.com/google/uuid"
)

type Agent interface {
	Activity()
	UpdateAgent()
	GetID() uuid.UUID
}
