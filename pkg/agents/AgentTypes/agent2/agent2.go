package agent2

import (
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
	"fmt"
)

type Agent2 struct {
	*baseAgent.BaseAgent
}

func (a2 *Agent2) UpdateAgent() {
	fmt.Println("Updating Agent2...")
}

func (a2 *Agent2) Activity() {
	fmt.Println("Agent2's Activity")
	a2.BaseAgent.Activity()
}

func GetAgent() baseAgent.Agent {
	return &Agent2{
		baseAgent.NewAgent("A2"),
	}
}
