package testUtils

import (
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
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

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration) *TestServer {
	m := make([]agent.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = agent.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	return &TestServer{
		BaseServer:       server.CreateServer(m, iterations, turns, maxDuration),
		TurnCounter:      0,
		IterationCounter: 0,
	}
}

func CreateTestTimeoutMessage(workLoad time.Duration) *TestTimeoutMessage {
	return &TestTimeoutMessage{
		message.BaseMessage{},
		workLoad,
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

func (t *TestServer) BroadcastMessage(msg message.IMessage[ITestBaseAgent]) {
	agMap := t.GetAgentMap()
	recipArr := make([]uuid.UUID, len(agMap))
	i := 0
	for id := range agMap {
		recipArr[i] = id
		i++
	}
	t.SendMessage(msg, recipArr)
}