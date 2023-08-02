package baseagent

import (
	"fmt"

	"github.com/google/uuid"
)

type BaseAgent struct {
	id   uuid.UUID
	name string
}

func (ba *BaseAgent) GetID() uuid.UUID {
	return ba.id
}

func (ba *BaseAgent) GetName() string {
	return ba.name
}

func (ba *BaseAgent) SetName(name string) {

	ba.name = name
}

func (ba *BaseAgent) UpdateAgent() {
	fmt.Println("UdateBaseAgent")
}

func (ba *BaseAgent) Activity() {
	fmt.Printf("id: %s\n", ba.GetID())
	fmt.Printf("name: %s\n", ba.GetName())
	fmt.Println("__________________________")
}

func NewAgent(name string) *BaseAgent {
	return &BaseAgent{
		id:   uuid.New(),
		name: name,
	}
}

func GetAgent() *BaseAgent {
	return &BaseAgent{
		id:   uuid.New(),
		name: "null",
	}
}
