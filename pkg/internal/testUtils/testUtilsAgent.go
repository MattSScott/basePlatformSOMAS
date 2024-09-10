package testUtils

import (
	"sync"
	"sync/atomic"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/google/uuid"
)

type ITestBaseAgent interface {
	agent.IAgent[ITestBaseAgent]
	NewTestMessage() TestMessage
	HandleTestMessage()
	ReceivedMessage() bool
	GetCounter() int32
	SetCounter(int32)
	GetGoal() int32
	SetGoal(int32)
	FinishedMessaging()
	NotifyAgentFinishedMessagingUnthreaded(*sync.WaitGroup, *uint32)
	GetAgentStoppedTalking() int
}

type TestServerFunctionsAgent struct {
	Counter        int32
	Goal           int32
	StoppedTalking int
	*agent.BaseAgent[ITestBaseAgent]
}

type IBadAgent interface {
	agent.IAgent[IBadAgent]
}

type AgentWithState struct {
	*agent.BaseAgent[ITestBaseAgent]
	State int
}

func (aws *AgentWithState) UpdateAgentInternalState() {
	aws.State += 1
}

func (ta *TestServerFunctionsAgent) FinishedMessaging() {
	ta.StoppedTalking++
	ta.NotifyAgentFinishedMessaging()
}

func (tba *TestServerFunctionsAgent) CreateTestMessage() TestMessage {
	return TestMessage{
		message.BaseMessage{},
		5,
	}
}

func (ag *TestServerFunctionsAgent) NotifyAgentFinishedMessagingUnthreaded(wg *sync.WaitGroup, counter *uint32) {
	defer wg.Done()
	ag.AgentStoppedTalking(ag.GetID())
	atomic.AddUint32(counter, 1)
}

func (ta TestServerFunctionsAgent) GetAgentStoppedTalking() int {
	return ta.StoppedTalking
}

func (ag *TestServerFunctionsAgent) HandleTestMessage() {
	newCounterValue := atomic.AddInt32(&ag.Counter, 1)
	if newCounterValue == atomic.LoadInt32(&ag.Goal) {
		ag.NotifyAgentFinishedMessaging()
	}
}

func (ag *TestServerFunctionsAgent) ReceivedMessage() bool {
	return ag.Counter == ag.Goal
}

func (ag *TestServerFunctionsAgent) UpdateAgentInternalState() {
	ag.Counter += 1
}

func NewTestAgent(serv agent.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestServerFunctionsAgent{
		BaseAgent:      agent.CreateBaseAgent(serv),
		Counter:        0,
		Goal:           0,
		StoppedTalking: 0,
	}
}

func (ta *TestServerFunctionsAgent) NewTestMessage() TestMessage {
	return TestMessage{
		ta.CreateBaseMessage(),
		5,
	}
}

func (ta *TestServerFunctionsAgent) GetCounter() int32 {
	return ta.Counter
}
func (ag *TestServerFunctionsAgent) RunSynchronousMessaging() {
	recipients := ag.ViewAgentIdSet()
	recipientArr := make([]uuid.UUID, len(recipients))
	i := 0
	for recip := range recipients {
		recipientArr[i] = recip
		i += 1
	}
	newMsg := ag.NewTestMessage()
	ag.SendSynchronousMessage(&newMsg, recipientArr)
}

func (tba *TestServerFunctionsAgent) SetCounter(count int32) {
	tba.Counter = count
}

func (tba *TestServerFunctionsAgent) SetGoal(goal int32) {
	tba.Goal = goal
}
func (tba *TestServerFunctionsAgent) GetGoal() int32 {
	return tba.Goal
}
