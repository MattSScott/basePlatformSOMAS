package agent2

import (
	baseUserAgent "basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
	"math/rand"

	"fmt"
)

type Agent2 struct {
	*baseUserAgent.AgentUser
	gender string
}

func setGender() string {
	if rand.Float32() < 0.5 {
		return "Male"
	}
	return "Female"

}

func (a1 *Agent2) Activity() {
	fmt.Println("Agent1's Activity")
	fmt.Printf("gender: %s\n", a1.gender)
	a1.BaseAgent.Activity()

}

func (a1 *Agent2) UpdateAgent() {
	fmt.Println("Updating Agent1...")
}

func GetAgent() baseAgent.Agent {
	return &Agent2{
		AgentUser: baseUserAgent.GetAgent("A2"),
		gender:    setGender(),
	}

}
