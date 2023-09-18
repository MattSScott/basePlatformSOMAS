package baseUserAgent

import (
	"fmt"

	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
)

type AgentUserInterface interface {
	baseAgent.Agent
	Activity1()
	Activity2()
}

type AgentUser struct {
	*baseAgent.BaseAgent
	name string
	food int
}

//add func for

func (a *AgentUser) Activity1() {
	fmt.Println("Agent's Activity1")

	fmt.Printf("name: %s\n", a.name)
	a.BaseAgent.Activity()
	a.SetFood()

}

func (a *AgentUser) Activity2() {
	fmt.Println("Agent's Activity2")

	// fmt.Printf("name: %s\n", a.name)
	// a.BaseAgent.Activity()
	// a.SetFood()

}

func GetAgent(name string) *AgentUser {
	return &AgentUser{
		BaseAgent: baseAgent.NewAgent(),
		name:      name,
		food:      0,
	}

}

func (u *AgentUser) GetName() string {
	return u.name
}

func (u *AgentUser) GetFood() int {
	return u.food
}
func (u *AgentUser) SetFood() {
	u.food++
}

func GetAgentDefault() baseAgent.Agent {
	return &AgentUser{
		BaseAgent: baseAgent.NewAgent(),
		name:      "BaseUser",
		food:      0,
	}
}
