package baseagent

import (
	"fmt"

	"github.com/google/uuid"
)

type BaseAgent struct {
	id uuid.UUID
}

func (ba *BaseAgent) GetID() uuid.UUID {
	return ba.id
}

// func (ba *BaseAgent) GetName() string {
// 	return ba.name
// }

// func (ba *BaseAgent) SetName(name string) {

// 	ba.name = name
// }

func (ba *BaseAgent) UpdateAgent() {
	fmt.Println("Updating BaseAgent...")
}

func (ba *BaseAgent) Activity() {
	fmt.Printf("id: %s\n", ba.GetID())
	fmt.Println("__________________________")
}

func NewAgent() *BaseAgent {
	return &BaseAgent{
		id: uuid.New(),
	}
}

func GetAgent() Agent {
	return &BaseAgent{
		id: uuid.New(),
	}
}
