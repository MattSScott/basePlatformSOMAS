package baseagent

import (
	"fmt"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
	"github.com/google/uuid"
)

type BaseAgent struct {
	id uuid.UUID
	msg string
	receivers []BaseAgent
	network []BaseAgent
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

func (ba *BaseAgent) GetMsg() string {
	return ba.msg
}

func (ba *BaseAgent) SetMsg(s string) {
	ba.msg = s
}

func (ba *BaseAgent) GetNet() []BaseAgent {
	return ba.network
}

func (ba *BaseAgent) SetNet(a []BaseAgent ) {
	ba.network = a 
}

func (ba *BaseAgent) GetRcv() []BaseAgent {
	return ba.receivers
}

func (ba *BaseAgent) SetRcv(a []BaseAgent) {
	ba.receivers = a
}

func (a *BaseAgent) GetMessage() message.Message[Agent] {
	return message.Message[Agent]{}
}
func (a *BaseAgent) HandleMessage(m message.Message[Agent]) message.Message[Agent] {
	return message.Message[Agent]{}
}