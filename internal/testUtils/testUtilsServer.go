package testUtils

import (
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type ITestServer interface {
	server.IServer[ITestBaseAgent]
}

type TestServer struct {
	*server.BaseServer[ITestBaseAgent]
	TurnCounter      int
	IterationCounter int
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int) *TestServer {
	serv := &TestServer{
		BaseServer:       server.CreateServer[ITestBaseAgent](iterations, turns, maxDuration, maxThreads),
		TurnCounter:      0,
		IterationCounter: 0,
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(NewTestAgent(serv))
	}
	return serv
}

func CreateTestTimeoutMessage(workLoad time.Duration) *TestTimeoutMessage {
	return &TestTimeoutMessage{
		message.BaseMessage{},
		workLoad,
	}
}

func CreateInfLoopMessage() *TestMessagingBandwidthLimiter {
	return &TestMessagingBandwidthLimiter{
		message.BaseMessage{},
	}
}

func NewTestMessage() *TestMessage {
	return &TestMessage{
		message.BaseMessage{},
		5,
	}
}

func (ts *TestServer) RunIteration() {
	ts.IterationCounter += 1
}

func (ts *TestServer) RunTurn() {
	ts.TurnCounter += 1
}

func (ts *TestServer) GetTurnCounter() int {
	return ts.TurnCounter
}

func (ts *TestServer) GetIterationCounter() int {
	return ts.IterationCounter
}

func SendNotifyMessages(agMap map[uuid.UUID]ITestBaseAgent, count *uint32, wg *sync.WaitGroup) {
	for _, ag := range agMap {
		wg.Add(1)
		go ag.NotifyAgentFinishedMessagingUnthreaded(wg, count)
	}
}
