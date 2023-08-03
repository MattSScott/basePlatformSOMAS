package baseUserAgent

import (
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
	"fmt"
)

type AgentUser struct {
	*baseAgent.BaseAgent
	name string
}

//add func for

func (a *AgentUser) Activity() {
	fmt.Println("Agent1's Activity")

	fmt.Printf("name: %s\n", a.name)
	a.BaseAgent.Activity()

}

func GetAgent(name string) *AgentUser {
	return &AgentUser{
		BaseAgent: baseAgent.NewAgent(),
		name:      name,
	}

}

func GetAgentDefault() baseAgent.Agent {
	return &AgentUser{
		BaseAgent: baseAgent.NewAgent(),
		name:      "BaseUser",
	}
}
