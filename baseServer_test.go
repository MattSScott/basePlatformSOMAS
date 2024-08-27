package basePlatformSOMAS_test

import (
	"sync"
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS"
	"github.com/google/uuid"
)

type ITestBaseAgent interface {
	basePlatformSOMAS.IAgent[ITestBaseAgent]
	NewTestMessage() TestMessage
	HandleTestMessage()
	ReceivedMessage() bool
	GetCounter() int
	SetCounter(int)
	GetGoal() int
	SetGoal(int)
}

type ITestServer interface {
	basePlatformSOMAS.IServer[ITestBaseAgent]
	basePlatformSOMAS.PrivateServerFields[ITestBaseAgent]
}

type TestAgent struct {
	*basePlatformSOMAS.BaseAgent[ITestBaseAgent]
	counter int
	goal    int
	mu      *sync.Mutex
}

type TestServer struct {
	*basePlatformSOMAS.BaseServer[ITestBaseAgent]
	roundCounter int
}

type TestMessage struct {
	basePlatformSOMAS.BaseMessage
	counter int
}

func (tm TestMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	ag.HandleTestMessage()
}

func (tba *TestAgent) CreateTestMessage() TestMessage {
	return TestMessage{
		basePlatformSOMAS.BaseMessage{},
		5,
	}
}

func (tba *TestAgent) SetCounter(count int) {
	tba.counter = count
}

func (tba *TestAgent) SetGoal(goal int) {
	tba.goal = goal
}
func (tba *TestAgent) GetGoal() int {
	return tba.goal
}
func NewTestMessage() TestMessage {
	return TestMessage{
		basePlatformSOMAS.BaseMessage{},
		5,
	}
}

func NewTestServer(generatorArray []basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], iterations, turns int, maxDuration time.Duration) *TestServer {
	return &TestServer{
		BaseServer:   basePlatformSOMAS.CreateServer(generatorArray, iterations, turns, maxDuration),
		roundCounter: 0,
	}
}

func NewTestAgent(serv basePlatformSOMAS.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	mu := &sync.Mutex{}
	return &TestAgent{
		BaseAgent: basePlatformSOMAS.CreateBaseAgent(serv),
		counter:   0,
		goal:      0,
		mu:        mu,
	}
}

func (ta2 *TestAgent) NewTestMessage() TestMessage {
	return TestMessage{
		basePlatformSOMAS.CreateBaseMessage(ta2.GetID()),
		5,
	}
}

func (ta2 *TestAgent) GetCounter() int {
	return ta2.counter
}
func (ag *TestAgent) RunSynchronousMessaging() {
	recipients := ag.ViewAgentIdSet()
	recipientArr := make([]uuid.UUID, len(recipients))
	i := 0
	for recip := range recipients {
		recipientArr[i] = recip
		i += 1
	}
	newMsg := ag.NewTestMessage()
	ag.SendSynchronousMessage(newMsg, recipientArr)
}

func (ts *TestServer) RunTurn() {
	ts.roundCounter += 1
}


func (ag *TestAgent) HandleTestMessage() {
	ag.mu.Lock()
	ag.counter += 1
	if ag.counter == ag.goal {
		ag.NotifyAgentFinishedMessaging()
	}
	ag.mu.Unlock()
}

func (ag *TestAgent) ReceivedMessage() bool {
	if ag.counter == ag.goal {
		return true
	} else {
		return false
	}

}

func (ag *TestAgent) UpdateAgentInternalState() {
	ag.counter += 1
}

func (server *TestServer) HandleTurn() {
	server.roundCounter += 1
}

func TestGenerateServer(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	if len(server.GetAgentMap()) != 3 {
		t.Error("len of agentmap is ", len(server.GetAgentMap()))
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)

	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)

	ag := NewTestAgent(server)
	ag.NotifyAgentFinishedMessaging()
	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != 3 {
		t.Error("Incorrect number of agents added to server", lenAgentMap)
	}
}
func TestHandlerInitialiser(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when handler not set")
		}
	}()
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	server.Initialise()
	server.Start()
}

func TestGenerateArrayFromMap(t *testing.T) {
	mapFound := make(map[uuid.UUID]int)
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agentArray := server.GenerateAgentArrayFromMap()
	for _, agent := range agentArray {
		_, exists := mapFound[agent.GetID()]
		if exists {
			t.Error("Duplicate")
		} else {
			mapFound[agent.GetID()] = 1
		}
	}
	for id := range server.GetAgentMap() {
		_, exists := mapFound[id]
		if !exists {
			t.Error("Value not in array")
		}
	}
}

func TestRunTurn(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	iterations := 1
	rounds := 2
	server := NewTestServer(m, iterations, rounds, time.Microsecond)
	server.SetRunHandler(server)
	server.Start()
	if server.roundCounter != (iterations * rounds) {
		t.Error("wrong number of rounds")
	}
}

func TestAgentRecievesMessage(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agent1 := NewTestAgent(server)
	testMessage := agent1.NewTestMessage()

	arrayReceivers := make([]uuid.UUID, 1)
	i := 0
	for id, ag := range server.GetAgentMap() {
		arrayReceivers[i] = id
		i += 1
		ag.SetGoal(1)
	}

	server.Initialise()
	go server.SendMessage(testMessage, arrayReceivers)
	server.EndAgentListeningSession()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error("Didn't Receive Message")
		}
	}
}

func TestWaitForMessagingToEnd(t *testing.T) {
	//This tests the function WaitForMessagingToEnd() which serves as a synchronisation procedure to ensure the main
	//thread waits for messaging to complete before moving the main thread on
	numberOfMessages := 10099
	const numAgents = 10
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id, ag := range agentMap {
		arrayOfIDs[i] = id
		i++
		ag.SetGoal(numberOfMessages *(numAgents))
	}

	for j := 0; j < numberOfMessages; j++ {
		for _,ag := range server.GetAgentMap() {
		msg := ag.NewTestMessage()
		//server.IncrementWaitGroup()
		go ag.SendMessage(msg, arrayOfIDs)
	}}
	server.EndAgentListeningSession()
	for _, ag := range agentMap {

		if !ag.ReceivedMessage() {
			t.Errorf("agent %s recieved %d messages, expected %d\n", ag.GetID(), ag.GetCounter(), ag.GetGoal())
		}
	}

}

func TestAddAgent(t *testing.T) {
	const numAgents = 10000
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agent := NewTestAgent(server)
	server.AddAgent(agent)
	agMap := server.GetAgentMap()

	if len(agMap) != numAgents+1 {
		t.Error("Adding Agent Failed")
	}
}

func TestRemoveAgent(t *testing.T) {
	const numAgents = 10000
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agMap := server.GetAgentMap()
	for _, ag := range agMap {
		server.RemoveAgent(ag)
	}

	if len(agMap) != 0 {
		t.Error("Removing Agent Failed")
	}
}

func TestNumIterationsInServer(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)

	if server.GetIterations() != 1 {
		t.Error("Incorrect number of iterations instantiated")
	}
}

func TestSendSynchronousMessase(t *testing.T) {

	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agent1 := NewTestAgent(server)
	testMessage := agent1.NewTestMessage()

	arrayReceivers := make([]uuid.UUID, 1)
	i := 0
	for id, ag := range server.GetAgentMap() {
		arrayReceivers[i] = id
		i += 1
		ag.SetGoal(1)
	}

	server.Initialise()
	server.SendSynchronousMessage(testMessage, arrayReceivers)
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error("Didn't Receive Message")
		}

	}
}

func TestSynchronousMessagingSession(t *testing.T) {
	numberAgents := 5
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numberAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	server.RunSynchronousMessagingSession()
	for _, ag := range server.GetAgentMap() {

		if ag.GetCounter() != numberAgents-1 {
			t.Error("All messages did not pass", ag.GetCounter())
		}
	}
}

func TestAccessAgentByID(t *testing.T) {
	numberAgents := 10
	randNum := 2347
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numberAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Microsecond)

	for _, ag := range server.GetAgentMap() {

		ag.SetCounter(randNum)

	}
	for id := range server.GetAgentMap() {
		if server.AccessAgentByID(id).GetCounter() != randNum {
			t.Error("Access Agent By ID is not working (incorrect struct value in test agent)")
		}

	}
}

func TestMessagePrint(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Microsecond)
	ag := NewTestAgent(server)
	newMsg := ag.NewTestMessage()
	newMsg.Print()
}
