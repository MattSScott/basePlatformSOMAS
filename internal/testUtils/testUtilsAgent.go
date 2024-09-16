package testUtils

import (
	"sync"
	"sync/atomic"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/google/uuid"
)

type ITestBaseAgent interface {
	agent.IAgent[ITestBaseAgent]
	CreateTestMessage() *TestMessage
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

func (ta *TestServerFunctionsAgent) UpdateAgentInternalState() {
	ta.Counter += 1
}

func (ta *TestServerFunctionsAgent) FinishedMessaging() {
	ta.StoppedTalking++
	ta.NotifyAgentFinishedMessaging()
}

func (ta *TestServerFunctionsAgent) CreateTestMessage() *TestMessage {
	return &TestMessage{
		ta.CreateBaseMessage(),
		5,
	}
}

func (ta *TestServerFunctionsAgent) NotifyAgentFinishedMessagingUnthreaded(wg *sync.WaitGroup, counter *uint32) {
	defer wg.Done()
	ta.AgentStoppedTalking(ta.GetID())
	atomic.AddUint32(counter, 1)
}

func (ta TestServerFunctionsAgent) GetAgentStoppedTalking() int {
	return ta.StoppedTalking
}

func (ta *TestServerFunctionsAgent) HandleTestMessage() {
	newCounterValue := atomic.AddInt32(&ta.Counter, 1)
	if newCounterValue == atomic.LoadInt32(&ta.Goal) {
		ta.NotifyAgentFinishedMessaging()
	}
}

func (ta *TestServerFunctionsAgent) ReceivedMessage() bool {
	return ta.Counter == ta.Goal
}

func NewTestAgent(serv agent.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestServerFunctionsAgent{
		BaseAgent:      agent.CreateBaseAgent(serv),
		Counter:        0,
		Goal:           0,
		StoppedTalking: 0,
	}
}

func (ta *TestServerFunctionsAgent) GetCounter() int32 {
	return ta.Counter
}
func (ta *TestServerFunctionsAgent) RunSynchronousMessaging() {
	recipients := ta.ViewAgentIdSet()
	recipientArr := make([]uuid.UUID, len(recipients))
	i := 0
	for recip := range recipients {
		recipientArr[i] = recip
		i += 1
	}
	newMsg := ta.CreateTestMessage()
	ta.SendSynchronousMessage(newMsg, recipientArr)
}

func (ta *TestServerFunctionsAgent) SetCounter(count int32) {
	ta.Counter = count
}

func (ta *TestServerFunctionsAgent) SetGoal(goal int32) {
	ta.Goal = goal
}
func (ta *TestServerFunctionsAgent) GetGoal() int32 {
	return ta.Goal
}
