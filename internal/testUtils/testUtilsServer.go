package testUtils

import (
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/google/uuid"
)

type ITestServer interface {
	server.IServer[ITestBaseAgent]
}

type TestTurnMethodPanics struct {
	*server.BaseServer[ITestBaseAgent]
}

type TestServer struct {
	*server.BaseServer[ITestBaseAgent]
	TurnCounter           int
	IterationStartCounter int
	IterationEndCounter   int
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int) *TestServer {
	serv := &TestServer{
		BaseServer:            server.CreateServer[ITestBaseAgent](iterations, turns, maxDuration, maxThreads),
		TurnCounter:           0,
		IterationStartCounter: 0,
		IterationEndCounter:   0,
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

func (ts *TestServer) RunTurn(turn, iteration int) {
	ts.TurnCounter += 1
}

func (ts *TestServer) RunStartOfIteration(iteration int) {
	ts.IterationStartCounter += 1
}

func (ts *TestServer) RunEndOfIteration(iteration int) {
	ts.IterationEndCounter += 1
}

func SendNotifyMessages(agMap map[uuid.UUID]ITestBaseAgent, count *uint32, wg *sync.WaitGroup) {
	for _, ag := range agMap {
		wg.Add(1)
		go ag.NotifyAgentFinishedMessagingUnthreaded(wg, count)
	}
}


