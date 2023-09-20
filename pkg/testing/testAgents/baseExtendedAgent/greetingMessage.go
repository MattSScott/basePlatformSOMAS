package baseExtendedAgent

import (
	"github.com/MattSScott/basePlatformSOMAS/pkg/main/messaging"
)

type IExtendedMessage interface {
	messaging.IMessage[IExtendedMessage]
	GetGreeting() string
}

type GreetingMessage struct {
	messaging.BaseMessage[IExtendedAgent]
	greeting string
}

func (em *GreetingMessage) GetGreeting() string {
	return em.greeting
}

func (em GreetingMessage) Accept(agent IExtendedAgent) {
	agent.HandleGreetingMessage(em)
}

func CreateGreetingMessage(sender IExtendedAgent, recipients []IExtendedAgent, content string) GreetingMessage {
	return GreetingMessage{
		BaseMessage: messaging.CreateMessage[IExtendedAgent](sender, recipients),
		greeting:    content,
	}
}
