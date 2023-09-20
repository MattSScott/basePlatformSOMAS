package baseExtendedAgent

import (
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/main/BaseAgent"
)

type ExtraMessagingFunctions interface {
	HandleGreetingMessage(GreetingMessage)
}

type IExtendedAgent interface {
	baseAgent.IAgent[IExtendedAgent]
	ExtraMessagingFunctions
	GetPhrase() string
}

type BaseExtendedAgent struct {
	*baseAgent.BaseAgent[IExtendedAgent]
	phrase string
}

func (ag *BaseExtendedAgent) GetPhrase() string {
	return ag.phrase
}

func GetAgent(phrase string) *BaseExtendedAgent {
	return &BaseExtendedAgent{
		BaseAgent: baseAgent.NewAgent[IExtendedAgent](),
		phrase:    phrase,
	}

}
