package main

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type IHelloWorldAgent interface {
	agent.IAgent[IHelloWorldAgent]
}

type IHelloWorldServer interface {
	server.IServer[IHelloWorldAgent]
}

type HelloWorldServer struct {
	*server.BaseServer[IHelloWorldAgent]
}

type HelloWorldAgent struct {
	*agent.BaseAgent[IHelloWorldAgent]
}

func CreateHelloWorldServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int) *HelloWorldServer {
	serv := &HelloWorldServer{
		BaseServer: server.CreateBaseServer[IHelloWorldAgent](iterations, turns, maxDuration, maxThreads),
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(CreateHelloWorldAgent(serv))
	}
	return serv
}

func CreateHelloWorldAgent(serv agent.IExposedServerFunctions[IHelloWorldAgent]) IHelloWorldAgent {
	return &HelloWorldAgent{
		BaseAgent: agent.CreateBaseAgent(serv),
	}
}

type HelloMessage struct {
	message.BaseMessage
}

type WorldMessage struct {
	message.BaseMessage
}

func (d *HelloMessage) InvokeMessageHandler(ag IHelloWorldAgent) {
	fmt.Print("Hello\n")
	msg := WorldMessage{ag.CreateBaseMessage()}
	ag.BroadcastMessage(&msg)
}

func (d *WorldMessage) InvokeMessageHandler(ag IHelloWorldAgent) {
	fmt.Print("World\n")
	// msg := HelloMessage{ag.CreateBaseMessage()}
	// ag.BroadcastMessage(&msg)
	ag.NotifyAgentFinishedMessaging()
}

func (serv *HelloWorldServer) RunTurn(i, j int) {
	fmt.Printf("Running turn %v,%v\n", i, j)
	k := 0
	for _, ag := range serv.GetAgentMap() {
		if k%2 == 0 {
			msg := HelloMessage{ag.CreateBaseMessage()}
			ag.BroadcastMessage(&msg)
		} else {
			msg := WorldMessage{ag.CreateBaseMessage()}
			ag.BroadcastMessage(&msg)
		}
	}
}

func (serv *HelloWorldServer) RunStartOfIteration(iteration int) {
	fmt.Printf("Starting iteration %v\n", iteration)
}

func (serv *HelloWorldServer) RunEndOfIteration(iteration int) {
	fmt.Printf("Ending iteration %v\n", iteration)
}

func main() {
	serv := CreateHelloWorldServer(10, 1, 1, time.Second, 10)
	serv.SetGameRunner(serv)
	serv.ReportMessagingDiagnostics()
	serv.Start()
}
