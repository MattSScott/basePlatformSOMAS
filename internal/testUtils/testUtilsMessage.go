package testUtils

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
)

type Message1 struct {
	message.BaseMessage
	MessageField1 int
}

type Message2 struct {
	message.BaseMessage
	MessageField2 int
}

type NullMessage struct {
	message.BaseMessage
}

type IExtendedAgent interface {
	agent.IAgent[IExtendedAgent]
	GetAgentField() int
	SetAgentField(int)
	GetMessage1() *Message1
	GetMessage2() *Message2
	GetNullMessage() *NullMessage
	HandleMessage1(msg Message1)
	HandleMessage2(msg Message2)
	HandleNullMessage(msg NullMessage)
}

type TestMessagingAgent struct {
	*agent.BaseAgent[IExtendedAgent]
	AgentField int
}

type TestMessage struct {
	message.BaseMessage
	Value int
}

type TestTimeoutMessage struct {
	message.BaseMessage
	Workload time.Duration
}

func NewExtendedAgent(serv agent.IExposedServerFunctions[IExtendedAgent]) IExtendedAgent {
	return &TestMessagingAgent{
		BaseAgent:  agent.CreateBaseAgent(serv),
		AgentField: 0,
	}
}

func (timeoutM TestTimeoutMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	ag.HandleTimeoutTestMessage(timeoutM)
}

func (tm TestMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	ag.HandleTestMessage()
}

func (agent *TestMessagingAgent) HandleMessage1(msg Message1) {
	agent.AgentField += msg.MessageField1
}

func (agent *TestMessagingAgent) HandleMessage2(msg Message2) {
	agent.AgentField += msg.MessageField2

}

func (agent *TestMessagingAgent) HandleNullMessage(msg NullMessage) {
	fmt.Println("handled!")
	agent.AgentField = 0
	fmt.Println(agent.AgentField)
}

func (m1 Message1) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage1(m1)
}

func (m2 Message2) InvokeMessageHandler(agent IExtendedAgent) {
	agent.HandleMessage2(m2)
}

func (nm NullMessage) InvokeMessageHandler(agent IExtendedAgent) {
	fmt.Println("Handling!")
	agent.HandleNullMessage(nm)
}

func (agent *TestMessagingAgent) GetAgentField() int {
	return agent.AgentField
}

func (agent *TestMessagingAgent) GetMessage1() *Message1 {
	return &Message1{
		agent.CreateBaseMessage(),
		5,
	}
}

func (agent *TestMessagingAgent) GetMessage2() *Message2 {
	return &Message2{
		agent.CreateBaseMessage(),
		10,
	}
}

func (agent *TestMessagingAgent) GetNullMessage() *NullMessage {
	return &NullMessage{
		agent.CreateBaseMessage(),
	}
}

func (agent *TestMessagingAgent) SetAgentField(num int) {
	agent.AgentField = num
}
