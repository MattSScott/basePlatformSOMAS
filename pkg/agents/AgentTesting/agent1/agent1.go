package agent1

import (
	baseUserAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	messaging "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
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

func (a1 *Agent1) GetMessage() messaging.Message[baseUserAgent.AgentUserInterface] {
	return messaging.Message[baseUserAgent.AgentUserInterface]{}
}

func (a1 *Agent1) HandleMessage(msg messaging.Message[baseUserAgent.AgentUserInterface]) messaging.Message[baseUserAgent.AgentUserInterface] {
	content := strings.ToLower(msg.GetContent())
	recip := []baseUserAgent.AgentUserInterface{msg.GetSender()}

	// TODO: very sloppy handling.
	// We should have the function just return content, so that the server can package it into a messaging object
	var sender baseUserAgent.AgentUserInterface = a1

	// using generics, we can now gain access to the methods from the interface
	msg.GetSender().Activity1()

	if content == "hello" {
		reply := messaging.CreateMessage(sender, "world", recip)
		return reply
	}

	return messaging.CreateMessage(sender, "sorry, I didn't quite get that", recip)
}
