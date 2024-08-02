package infra

import (
	"github.com/google/uuid"
)

type BaseAgent struct {
	id uuid.UUID
}

func (ba *BaseAgent) GetID() uuid.UUID {
	return ba.id
}

func CreateBaseAgent() *BaseAgent {
	return &BaseAgent{
		id: uuid.New(),
	}
}
