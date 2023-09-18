package agent1

import (
	"fmt"
	"strings"

	baseUserAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	messaging "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
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
	return messaging.CreateMessage(a1, "hello", a1.GetNetworkForMessaging())
}

func (a1 *Agent1) HandleMessage(msg messaging.Message) {
	content := strings.ToLower(msg.GetContent())

	if content == "wello" {
		fmt.Println("horld")
	}

}
