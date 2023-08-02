package agent1

import (
	"fmt"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
)

type Agent1 struct {
	*baseAgent.BaseAgent
	age int
}

func (a1 *Agent1) Activity() {
	fmt.Println("Agent1's Activity")
	fmt.Printf("age: %d\n", a1.age)
	a1.BaseAgent.Activity()

}

func GetAgent() *Agent1 {
	return &Agent1{
		BaseAgent: baseAgent.NewAgent("A1"),
		age:       0,
	}

}
