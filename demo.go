package main

import (
	"fmt"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type IDemoAgent interface {
	agent.IAgent[IDemoAgent]
}

type IDemoServer interface {
	server.IServer[IDemoAgent]
}

type DemoServer struct {
	*server.BaseServer[IDemoAgent]
	TurnCounter  int
	RoundCounter int
}

type DemoAgent struct {
	*agent.BaseAgent[IDemoAgent]
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration) *DemoServer {
	m := make([]agent.AgentGeneratorCountPair[IDemoAgent], 1)
	m[0] = agent.MakeAgentGeneratorCountPair(NewDemoAgent, numAgents)
	return &DemoServer{
		BaseServer: server.CreateServer(m, iterations, turns, maxDuration),
	}
}

func NewDemoAgent(serv agent.IExposedServerFunctions[IDemoAgent]) IDemoAgent {
	return &DemoAgent{
		BaseAgent: agent.CreateBaseAgent(serv),
	}
}

type DemoMessage1 struct {
	message.BaseMessage
}

type DemoMessage2 struct {
	message.BaseMessage
}

func (d *DemoMessage1) InvokeMessageHandler(ag IDemoAgent) {
	fmt.Print("Hello World\n")
	agSet := ag.ViewAgentIdSet()
	arrayRec := make([]uuid.UUID, len(agSet) - 1)
	i := 0
	for id := range agSet {
		if id != ag.GetID() {
			arrayRec[i] = id
			//fmt.Println(id)
			i++
		}
	}
	msg := DemoMessage2{ag.CreateBaseMessage()}
	time.Sleep(1 * time.Millisecond)
	ag.SendMessage(&msg, arrayRec)
}

func (d *DemoMessage2) InvokeMessageHandler(ag IDemoAgent) {
	fmt.Print("Wello Horld\n")
	agSet := ag.ViewAgentIdSet()

	arrayRec := make([]uuid.UUID, len(agSet) - 1)

	i := 0
	for id := range agSet {
		if id != ag.GetID() {
			arrayRec[i] = id
			i++
		}
	}
	msg := DemoMessage1{ag.CreateBaseMessage()}
	time.Sleep(1 * time.Millisecond)
	ag.SendMessage(&msg, arrayRec)
}

func (serv *DemoServer) RunTurn() {
	agMap := serv.GetAgentMap()
	recArray := make([]uuid.UUID, len(agMap))
	z := 0
	for id := range serv.GetAgentMap() {
		recArray[z] = id
		z++
	}
	i := 0
	for _, ag := range serv.GetAgentMap() {
		if i % 2== 0 {
			msg := DemoMessage1{ag.CreateBaseMessage()}
			ag.SendMessage(&msg, recArray)
		} else {
			msg := DemoMessage2{ag.CreateBaseMessage()}
			ag.SendMessage(&msg, recArray)
		}
		i++
	}
}

func main() {
	fmt.Println("Running Demo")
	numAgents := 2
	numTurns := 2
	numIterations := 2
	timeout := time.Millisecond
	serv := GenerateTestServer(numAgents, numIterations, numTurns, timeout)
	serv.SetGameRunner(serv)
	serv.Start()
}
