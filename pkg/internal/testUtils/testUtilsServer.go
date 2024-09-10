package testUtils

import (
	"fmt"
	"sync"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type ITestServer interface {
	server.IServer[ITestBaseAgent]
	server.PrivateServerFields[ITestBaseAgent]
}

type TestServer struct {
	*server.BaseServer[ITestBaseAgent]
	TurnCounter  int
	RoundCounter int
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration) *TestServer {
	m := make([]agent.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = agent.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	return &TestServer{
		BaseServer:   server.CreateServer(m, iterations, turns, maxDuration),
		TurnCounter:  0,
		RoundCounter: 0,
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

func NewTestMessage() TestMessage {
	return TestMessage{
		message.BaseMessage{},
		5,
	}
}

func (ts *TestServer) RunTurn() {
	ts.TurnCounter += 1
}

func (ts *TestServer) RunRound() {
	ts.RoundCounter += 1
}

func (server *TestServer) HandleTurn() {
	server.TurnCounter += 1
}

func (tba *TestServer) InfMessageSend(newMsg InfiniteLoopMessage, receiver []uuid.UUID, done chan struct{}) {
	go tba.SendMessage(&newMsg, receiver)
	fmt.Println(tba.EndAgentListeningSession())
	done <- struct{}{}
}

func (tba *TestServer) GetTurnCounter() int {
	return tba.TurnCounter
}

func (tba *TestServer) GetRoundCounter() int {
	return tba.RoundCounter
}

func SendNotifyMessages(agMap map[uuid.UUID]ITestBaseAgent, count *uint32, iter int, wg *sync.WaitGroup) {
	for _, ag := range agMap {
		for i := 0; i < iter; i++ {
			fmt.Println("running")
			wg.Add(1)
			go ag.NotifyAgentFinishedMessagingUnthreaded(wg, count)
		}
	}
}

type RunHandler struct {
	Iters int
	Turns int
}

func (r *RunHandler) RunRound() {
	r.Iters++
}
func (r *RunHandler) RunTurn() {
	r.Turns++
}
