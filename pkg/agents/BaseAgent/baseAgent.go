package baseagent

import (
	"fmt"

	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
	"github.com/google/uuid"
)

type BaseAgent struct {
	id      uuid.UUID
	network []BaseAgent
}

func (ba *BaseAgent) GetID() uuid.UUID {
	return ba.id
}

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

func (ba *BaseAgent) GetNetwork() []BaseAgent {
	return ba.network
}

func (ba *BaseAgent) GetNetworkForMessaging() []message.Messaging {
	messengerArray := make([]message.Messaging, len(ba.network))
	for i := range ba.network {
		messengerArray[i] = &ba.network[i]
	}
	return messengerArray
}

func (ba *BaseAgent) SetNetwork(newNetwork []BaseAgent) {
	ba.network = newNetwork
}

func (ba *BaseAgent) GetMessage() message.Message {
	return message.CreateMessage(ba, "", ba.GetNetworkForMessaging())
}
func (a *BaseAgent) HandleMessage(m message.Message) {
	fmt.Println("message received in baseAgent")
}
