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
	HandleHelloMessage(HelloMessage)
	HandleWorldMessage(WorldMessage)
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

func (hwa *HelloWorldAgent) HandleHelloMessage(msg HelloMessage) {
	fmt.Printf("%s said: 'Hello'\n", hwa.GetID())
	response := WorldMessage{hwa.CreateBaseMessage()}
	hwa.SendMessage(&response, msg.Sender)
}

func (hwa *HelloWorldAgent) HandleWorldMessage(msg WorldMessage) {
	fmt.Printf("%s responded: 'World'\n", hwa.GetID())
	hwa.NotifyAgentFinishedMessaging()
}

func CreateHelloWorldServer(numAgents, iterations, turns int, maxDuration time.Duration, maxThreads int) *HelloWorldServer {
	serv := &HelloWorldServer{
		BaseServer: server.CreateBaseServer[IHelloWorldAgent](iterations, turns, maxDuration, maxThreads),
	}
	for i := 0; i < numAgents; i++ {
		serv.AddAgent(CreateHelloWorldAgent(serv))
	}
	serv.SetGameRunner(serv)
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

func (d HelloMessage) InvokeMessageHandler(ag IHelloWorldAgent) {
	ag.HandleHelloMessage(d)
}

func (d WorldMessage) InvokeMessageHandler(ag IHelloWorldAgent) {
	ag.HandleWorldMessage(d)
}

func (serv *HelloWorldServer) RunTurn(i, j int) {
	fmt.Printf("Running iteration %v, turn %v\n", i+1, j+1)
	for _, ag := range serv.GetAgentMap() {
		msg := HelloMessage{ag.CreateBaseMessage()}
		ag.BroadcastMessage(&msg)
	}
}

func (serv *HelloWorldServer) RunStartOfIteration(iteration int) {
	fmt.Printf("Starting iteration %v\n", iteration+1)
	fmt.Println()
}

func (serv *HelloWorldServer) RunEndOfIteration(iteration int) {
	fmt.Println()
	fmt.Printf("Ending iteration %v\n", iteration+1)
}

func main() {
	serv := CreateHelloWorldServer(4, 1, 1, time.Millisecond, 100)
	serv.ReportMessagingDiagnostics()
	serv.Start()
}
