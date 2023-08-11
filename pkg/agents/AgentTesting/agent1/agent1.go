package agent1

import (
	baseUserAgent "basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
	messaging "basePlatformSOMAS/pkg/messaging"
	"fmt"
	"strings"
)

type Agent1 struct {
	*baseUserAgent.AgentUser
	age int
}

func (a1 *Agent1) Activity() {

	fmt.Println("Agent1's Activity")
	fmt.Printf("age: %d\n", a1.age)
	a1.AgentUser.Activity()

}

func (a1 *Agent1) UpdateAgent() {
	fmt.Println("Updating Agent1...")
	a1.age++
}

func (a1 *Agent1) GetAge() int {
	return a1.age
}

func GetAgent() baseUserAgent.AgentUserInterface {
	return &Agent1{
		AgentUser: baseUserAgent.GetAgent("A1"),
		age:       0,
	}
}

func (a1 *Agent1) GetMessage() messaging.Message {
	return messaging.Message{}
}

func (a1 *Agent1) HandleMessage(msg messaging.Message) messaging.Message {
	content := strings.ToLower(msg.GetContent())
	recip := []baseAgent.Agent{msg.GetSender()}

	if content == "hello" {
		reply := messaging.CreateMessage(a1, "world", recip)
		return reply
	}

	return messaging.CreateMessage(a1, "sorry, I didn't quite get that", recip)
}
