package main

import (
	"fmt"
	"sync/atomic"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/pkg/server"
	"github.com/google/uuid"
)

type IDemoAgent interface {
	agent.IAgent[IDemoAgent]
	ICounterFunctions
	SetICounterFunctions(ICounterFunctions)
}

type ICounterFunctions interface {
	IncrementCounter()
}

type IDemoServer interface {
	server.IServer[IDemoAgent]
	ICounterFunctions
	PrintCounter()
}

type DemoServer struct {
	ThreadCounter int32
	*server.BaseServer[IDemoAgent]
}

type DemoAgent struct {
	*agent.BaseAgent[IDemoAgent]
	ICounterFunctions
}

type DemoMessage1 struct {
	message.BaseMessage
}

type DemoMessage2 struct {
	message.BaseMessage
}

func GenerateTestServer(numAgents, iterations, turns int, maxDuration time.Duration) *DemoServer {
	m := make([]agent.AgentGeneratorCountPair[IDemoAgent], 1)
	m[0] = agent.MakeAgentGeneratorCountPair(NewDemoAgent, numAgents)
	server := DemoServer{
		BaseServer:    server.CreateServer(m, iterations, turns, maxDuration, 1000),
		ThreadCounter: 0,
	}
	for _, ag := range server.GetAgentMap() {
		ag.SetICounterFunctions(&server)
	}
	return &server
}

func NewDemoAgent(serv agent.IExposedServerFunctions[IDemoAgent]) IDemoAgent {
	return &DemoAgent{
		BaseAgent: agent.CreateBaseAgent(serv),
	}
}

func (a *DemoAgent) SetICounterFunctions(functions ICounterFunctions) {
	a.ICounterFunctions = functions
}

func (server *DemoServer) IncrementCounter() {
	atomic.AddInt32(&server.ThreadCounter, 1)
}

func (server *DemoServer) PrintCounter() {
	fmt.Println(server.ThreadCounter)
}

func (d *DemoMessage1) InvokeMessageHandler(ag IDemoAgent) {
	agSet := ag.ViewAgentIdSet()
	arrayRec := make([]uuid.UUID, len(agSet)-1)
	i := 0
	for id := range agSet {
		if id != ag.GetID() {
			arrayRec[i] = id
			i++
		}
	}
	msg := DemoMessage2{ag.CreateBaseMessage()}
	ag.SendMessage(&msg, arrayRec)
	ag.IncrementCounter()
	ag.NotifyAgentFinishedMessaging()
}

func (d *DemoMessage2) InvokeMessageHandler(ag IDemoAgent) {
	agSet := ag.ViewAgentIdSet()
	arrayRec := make([]uuid.UUID, len(agSet)-1)
	i := 0
	for id := range agSet {
		if id != ag.GetID() {
			arrayRec[i] = id
			i++
		}
	}
	msg := DemoMessage1{ag.CreateBaseMessage()}
	ag.IncrementCounter()
	ag.SendMessage(&msg, arrayRec)
	ag.NotifyAgentFinishedMessaging()
}

func (serv *DemoServer) RunTurn() {
	agMap := serv.GetAgentMap()
	lenAgMap := len(agMap)
	expectedNumMessages := lenAgMap * lenAgMap
	fmt.Println("Expecting number of messages:", expectedNumMessages)
	recArray := make([]uuid.UUID, len(agMap))
	z := 0
	for id := range serv.GetAgentMap() {
		recArray[z] = id
		z++
	}
	for _, ag := range serv.GetAgentMap() {
		msg1 := DemoMessage1{ag.CreateBaseMessage()}
		fmt.Println("sending message")
		ag.SendMessage(&msg1, recArray)
	}
}

func main() {
	fmt.Println("Running Demo")
	numAgents := 4
	numTurns := 1
	numIterations := 1
	timeout := time.Second
	serv := GenerateTestServer(numAgents, numIterations, numTurns, timeout)
	serv.SetGameRunner(serv)
	serv.Start()
	serv.PrintCounter()
	time.Sleep(time.Second)
}
