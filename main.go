package main

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type ITestAgent interface {
	agent.IAgent[ITestAgent]
}

type TestAgent struct {
	*agent.BaseAgent[ITestAgent]
}

type ITestServer interface {
	server.IServer[ITestAgent]
}

type TestServer struct {
	*server.BaseServer[ITestAgent]
}

func NewTestAgent(serv *TestServer) ITestAgent {
	return &TestAgent{
		BaseAgent: agent.CreateBaseAgent(serv)}
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int) *TestServer {
	serv := &TestServer{
		BaseServer: server.CreateServer[ITestAgent](iterations, turns, maxDuration, maxThreads),
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(NewTestAgent(serv))
	}
	return serv
}

func (s *TestServer) RunTurn(i, j int) {
	z := 0
	for _, ag := range s.GetAgentMap() {
		if !(z == 2 || z == 78 || z == 10) {
			ag.NotifyAgentFinishedMessaging()
		}
		z++
	}
	time.Sleep(10 * time.Millisecond)
}

func (s *TestServer) RunStartOfIteration(i int) {}
func (s *TestServer) RunEndOfIteration(i int)   {}

func main() {
	server := GenerateTestServer(100, 1, 1, 100*time.Millisecond, 100)
	server.SetGameRunner(server)
	fmt.Printf("Starting\n")
	server.Start()

}
