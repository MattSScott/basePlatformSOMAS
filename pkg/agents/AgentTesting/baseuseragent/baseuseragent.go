package baseUserAgent

import (
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
	"fmt"
)

type AgentUserInterface interface {
	baseAgent.Agent
	Activity1()
	Activity2()
	GetMessage() message.Message[AgentUserInterface]
	HandleMessage(message.Message[AgentUserInterface]) message.Message[AgentUserInterface]
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

func (a *AgentUser) GetMessage() message.Message[AgentUserInterface] {
	return message.Message[AgentUserInterface]{}
}
func (a *AgentUser) HandleMessage(m message.Message[AgentUserInterface]) message.Message[AgentUserInterface] {
	return message.Message[AgentUserInterface]{}
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
