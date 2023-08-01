package agent1

import (
	"fmt"
	baseAgent "somas_base_platform/pkg/agents/BaseAgent"
)

type Agent1 struct {
	*baseAgent.BaseAgent
	age int
	//Energy int
}

// func ReduceEnergy(a *Agent1) {
// 	a.Energy -= 10
// }

func (a1 *Agent1) Activity() {
	fmt.Println("Agent1's Activity")
	a1.BaseAgent.Activity()
	fmt.Printf("age: %d\n", a1.age)

	// fmt.Printf("id: %s\n", a1.GetID())
	// a1.SetName("A1")
	// fmt.Printf("name: %s\n", a1.GetName())
	// fmt.Printf("__________________________\n")
}

func GetAgent() *Agent1 { //why this function returns the interface?
	ag := &Agent1{
		BaseAgent: baseAgent.GetAgent(),
		age:       0,
	}

	ag.SetName("A1")

	return ag
}

//func (a *Agent1) ReduceEnergy() {
//	a.Energy -= 10
//}

// func (a *Agent1) UpdateAgent() {
// 	a.ReduceEnergy()
// 	a.SayHi()
// }

// func (a *Agent1) SendMessage() {
// 	fmt.Println("world")
// }
