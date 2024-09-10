package utils

import (
	"sync"
	"sync/atomic"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
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
	NotifyAgentFinishedMessagingUnthreaded(*sync.WaitGroup, *uint32)
}

type IBadAgent interface {
	agent.IAgent[IBadAgent]
}

type ITestServer interface {
	server.IServer[ITestBaseAgent]
	server.PrivateServerFields[ITestBaseAgent]
}

type TestAgent struct {
	counter int32
	goal    int32
	*agent.BaseAgent[ITestBaseAgent]
}

type TestServer struct {
	*server.BaseServer[ITestBaseAgent]
	turnCounter  int
	roundCounter int
}

type TestMessage struct {
	message.BaseMessage
	counter int
}

type InfiniteLoopMessage struct {
	message.BaseMessage
}

func (infM InfiniteLoopMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	InfLoop()
}

func (infM InfiniteLoopMessage) InvokeSyncMessageHandler(ag ITestBaseAgent) {
	InfLoop()
}

func (tm TestMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	ag.HandleTestMessage()
}

func (tba *TestAgent) CreateTestMessage() TestMessage {
	return TestMessage{
		message.BaseMessage{},
		5,
	}
}

func (ag *TestAgent) NotifyAgentFinishedMessagingUnthreaded(wg *sync.WaitGroup, counter *uint32) {
	defer wg.Done()
	ag.AgentStoppedTalking(ag.GetID())
	atomic.AddUint32(counter, 1)
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration) *TestServer {
	m := make([]agent.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = agent.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	return &TestServer{
		BaseServer:   server.CreateServer(m, iterations, turns, maxDuration),
		turnCounter:  0,
		roundCounter: 0,
	}
}

func CreateInfiniteLoopMessage() InfiniteLoopMessage {
	return InfiniteLoopMessage{
		message.BaseMessage{},
	}
}

func InfLoop() {
	i := 0
	for {
		if i == -1 {
			return
		}
		time.Sleep(1000 * time.Millisecond)
	}
}
func (tba *TestAgent) SetCounter(count int32) {
	tba.counter = count
}

func (tba *TestAgent) SetGoal(goal int32) {
	tba.goal = goal
}
func (tba *TestAgent) GetGoal() int32 {
	return tba.goal
}
func NewTestMessage() TestMessage {
	return TestMessage{
		message.BaseMessage{},
		5,
	}
}

func NewTestAgent(serv agent.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestAgent{
		BaseAgent: agent.CreateBaseAgent(serv),
		counter:   0,
		goal:      0,
	}
}

func (ta *TestAgent) NewTestMessage() TestMessage {
	return TestMessage{
		ta.CreateBaseMessage(),
		5,
	}
}

func (ta *TestAgent) GetCounter() int32 {
	return ta.counter
}
func (ag *TestAgent) RunSynchronousMessaging() {
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

func (ts *TestServer) RunTurn() {
	ts.turnCounter += 1
}

func (ts *TestServer) RunRound() {
	ts.roundCounter += 1
}

func (ag *TestAgent) HandleTestMessage() {
	newCounterValue := atomic.AddInt32(&ag.counter, 1)
	if newCounterValue == atomic.LoadInt32(&ag.goal) {
		ag.NotifyAgentFinishedMessaging()
	}
}

func (ag *TestAgent) ReceivedMessage() bool {
	return ag.counter == ag.goal
}

func (ag *TestAgent) UpdateAgentInternalState() {
	ag.counter += 1
}

func (server *TestServer) HandleTurn() {
	server.turnCounter += 1
}
