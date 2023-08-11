package agent1

import (
	baseUserAgent "basePlatformSOMAS/pkg/agents/AgentTesting/baseuseragent"
	messaging "basePlatformSOMAS/pkg/messaging" 

	"fmt"
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



func (a1 *Agent1) HandleMsg[T baseUserAgent.AgentUserInterface](letter messaging.LetterI[baseUserAgent.AgentUserInterface]){
	msg := letter.GetMsg()
	if msg == "Hello"{
		a1.SetMsg("world!")
		
		s:=letter.GetSnd()
		snd:= [1]T{s}

		a1.SetRcv(snd)
	}


}