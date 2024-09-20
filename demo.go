package main

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type IACSOSAgent interface {
	agent.IAgent[IACSOSAgent]
	Talk()
	HandleHelloMessage(HelloMessage)
	HandleWelloMessage(WelloMessage)
}

type IACSOSServer interface {
	server.IServer[IACSOSAgent]
}

type ACSOSServer struct {
	*server.BaseServer[IACSOSAgent]
}

type HelloAgent struct {
	*agent.BaseAgent[IACSOSAgent]
}

func GenerateHelloAgent(serv agent.IExposedServerFunctions[IACSOSAgent]) IACSOSAgent {
	return &HelloAgent{
		agent.CreateBaseAgent(serv),
	}
}

type WelloAgent struct {
	*agent.BaseAgent[IACSOSAgent]
}

func GenerateWelloAgent(serv agent.IExposedServerFunctions[IACSOSAgent]) IACSOSAgent {
	return &WelloAgent{
		agent.CreateBaseAgent(serv),
	}
}

type HelloMessage struct {
	message.BaseMessage
}

func (hm HelloMessage) InvokeMessageHandler(ag IACSOSAgent) {
	ag.HandleHelloMessage(hm)
}

type WelloMessage struct {
	message.BaseMessage
}

func (wm WelloMessage) InvokeMessageHandler(ag IACSOSAgent) {
	ag.HandleWelloMessage(wm)
}

func (ha *HelloAgent) HandleHelloMessage(msg HelloMessage) {
}

func (ha *HelloAgent) HandleWelloMessage(msg WelloMessage) {
	fmt.Println("Horld")
	ha.NotifyAgentFinishedMessaging()
}

func (ha *HelloAgent) Talk() {
	fmt.Println("Hello")
	ha.BroadcastMessage(&HelloMessage{ha.CreateBaseMessage()})
	ha.NotifyAgentFinishedMessaging()
}

func (ha *WelloAgent) Talk() {
	fmt.Println("Wello")
	ha.BroadcastMessage(&WelloMessage{ha.CreateBaseMessage()})
	ha.NotifyAgentFinishedMessaging()
}

func (wa *WelloAgent) HandleHelloMessage(msg HelloMessage) {
	fmt.Println("World")
	recip := []uuid.UUID{msg.GetSender()}
	wa.SendMessage(&WelloMessage{wa.CreateBaseMessage()}, recip)
	wa.NotifyAgentFinishedMessaging()
}

func (ha *WelloAgent) HandleWelloMessage(msg WelloMessage) {
}

func GenerateACSOSServer(numAgents, iterations, turns int, maxDuration time.Duration) *ACSOSServer {
	m := make([]agent.AgentGeneratorCountPair[IACSOSAgent], 2)
	m[0] = agent.MakeAgentGeneratorCountPair(GenerateHelloAgent, numAgents)
	m[1] = agent.MakeAgentGeneratorCountPair(GenerateWelloAgent, numAgents)

	return &ACSOSServer{
		BaseServer: server.CreateServer(m, iterations, turns, maxDuration, 1000),
	}

}

func (serv *ACSOSServer) RunTurn() {
	for _, ag := range serv.GetAgentMap() {
		ag.Talk()
	}
}

func main() {
	fmt.Println("Running Demo")
	numAgents := 10
	numTurns := 1
	numIterations := 1
	timeout := time.Millisecond
	serv := GenerateACSOSServer(numAgents, numIterations, numTurns, timeout)
	serv.SetGameRunner(serv)
	serv.Start()
}
