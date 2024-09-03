package basePlatformSOMAS_test

import (
	"fmt"
	"sync/atomic"
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
	GetCounter() int64
	SetCounter(int64)
	GetGoal() int64
	SetGoal(int64)
}

type IBadAgent interface {
	basePlatformSOMAS.IAgent[IBadAgent]
}

type ITestServer interface {
	basePlatformSOMAS.IServer[ITestBaseAgent]
	basePlatformSOMAS.PrivateServerFields[ITestBaseAgent]
}

type TestAgent struct {
	*basePlatformSOMAS.BaseAgent[ITestBaseAgent]
	counter int64
	goal    int64
}

type TestServer struct {
	*basePlatformSOMAS.BaseServer[ITestBaseAgent]
	turnCounter  int
	roundCounter int
}

type TestMessage struct {
	basePlatformSOMAS.BaseMessage
	counter int
}

type InfiniteLoopMessage struct {
	basePlatformSOMAS.BaseMessage
}

func (infM InfiniteLoopMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	InfLoop()
}

func (infM InfiniteLoopMessage) InvokeSyncMessageHandler(ag ITestBaseAgent) {
	InfLoop()
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

func CreateInfiniteLoopMessage() InfiniteLoopMessage {
	return InfiniteLoopMessage{
		basePlatformSOMAS.BaseMessage{},
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
func (tba *TestAgent) SetCounter(count int64) {
	tba.counter = count
}

func (tba *TestAgent) SetGoal(goal int64) {
	tba.goal = goal
}
func (tba *TestAgent) GetGoal() int64 {
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
		turnCounter:  0,
		roundCounter: 0,
	}
}

func NewTestAgent(serv basePlatformSOMAS.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestAgent{
		BaseAgent: basePlatformSOMAS.CreateBaseAgent(serv),
		counter:   0,
		goal:      0,
	}
}

func (ta *TestAgent) NewTestMessage() TestMessage {
	return TestMessage{
		ta.CreateBaseMessage(),
		5,
	}
}

func (ta *TestAgent) GetCounter() int64 {
	return ta.counter
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
	ts.turnCounter += 1
}

func (ts *TestServer) RunRound() {
	ts.roundCounter += 1
}

func (ag *TestAgent) HandleTestMessage() {
	newCounterValue := atomic.AddInt64(&ag.counter, 1)
	if newCounterValue == atomic.LoadInt64(&ag.goal) {
		ag.NotifyAgentFinishedMessaging()
	}
}

func (ag *TestAgent) ReceivedMessage() bool {
	return ag.counter == ag.goal
}

func (ag *TestAgent) UpdateAgentInternalState() {
	ag.counter += 1
}

func (server *TestServer) HandleTurn() {
	server.turnCounter += 1
}

func TestGenerateServer(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	numAgents := 3
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != numAgents {
		t.Error(lenAgentMap, "agents initialised,expected", numAgents)
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	numAgents := 3
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)

	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)

	ag := NewTestAgent(server)
	server.EndAgentListeningSession()
	go ag.NotifyAgentFinishedMessaging()

	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != numAgents {
		t.Error("Incorrect number of agents added to server,got", lenAgentMap, "expected", numAgents)
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
	numAgents := 10
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	agentArray := server.GenerateAgentArrayFromMap()
	for _, agent := range agentArray {
		_, exists := mapFound[agent.GetID()]
		if exists {
			t.Error("Duplicate of", agent.GetID(), "found in output array")
		} else {
			mapFound[agent.GetID()] = 1
		}
	}
	for id := range server.GetAgentMap() {
		_, exists := mapFound[id]
		if !exists {
			t.Error(id, "in agentMap but not in array")
		}
	}
}

func TestRunTurn(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	iterations := 1
	rounds := 1
	server := NewTestServer(m, iterations, rounds, time.Millisecond)
	server.SetRunHandler(server)
	server.Start()
	if server.turnCounter != (iterations * rounds) {
		t.Error("wrong number of iterations executed", server.turnCounter, "expected", iterations*rounds)
	}
}

func TestRunRound(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	rounds := 1
	server := NewTestServer(m, 1, 1, time.Millisecond)
	server.SetRunHandler(server)
	for i := 0; i < rounds; i++ {
		server.RunRound()
	}
	if server.roundCounter != rounds {
		t.Error("wrong number of rounds executed", server.roundCounter, "expected", rounds)
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
	_ = server.EndAgentListeningSession()
	for _, ag := range server.GetAgentMap() {
		fmt.Println()
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
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
		ag.SetGoal(int64(numberOfMessages * numAgents))
	}

	for j := 0; j < numberOfMessages; j++ {
		for _, ag := range server.GetAgentMap() {
			msg := ag.NewTestMessage()
			go ag.SendMessage(msg, arrayOfIDs)
		}
	}
	a := server.EndAgentListeningSession()
	if !a {
		t.Error("Messaging ended early")
	}

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
	lenAgMap := len(agMap)
	expectedNumAgents := numAgents + 1
	if lenAgMap != expectedNumAgents {
		t.Error("Removing Agents Failed,expected number of agents:", expectedNumAgents, ",got:", lenAgMap)
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
	lenAgMap := len(agMap)
	if lenAgMap != 0 {
		t.Error("Removing Agents Failed,expected number of agents: 0,got:", lenAgMap)
	}
}

func TestNumIterationsInServer(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	iterations := 1
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
	getIterationsValue := server.GetIterations()
	if getIterationsValue != iterations {
		t.Error("Incorrect number of iterations instantiated, expected:", iterations, "got:", getIterationsValue)
	}
}

func TestSendSynchronousMessage(t *testing.T) {
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
	server.EndAsyncMessaging()
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
	server.EndAsyncMessaging()
	server.RunSynchronousMessagingSession()
	for _, ag := range server.GetAgentMap() {

		if ag.GetCounter() != int64(numberAgents-1) {
			t.Error("All messages did not pass, got:", ag.GetCounter(), "expected:", numberAgents-1)
		}
	}
}

func TestAccessAgentByID(t *testing.T) {
	numberAgents := 10
	var randNum int64 = 2357
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, numberAgents)
	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)

	for _, ag := range server.GetAgentMap() {

		ag.SetCounter(randNum)

	}
	for id := range server.GetAgentMap() {
		accessedAgentID := server.AccessAgentByID(id).GetCounter()
		if accessedAgentID != randNum {
			t.Error("Access Agent By ID is not working (incorrect struct value in test agent),expected:", randNum, "got:", accessedAgentID)
		}

	}
}

func TestMessagePrint(t *testing.T) {
	ag := NewTestAgent(nil)
	msg := ag.CreateBaseMessage()
	msg.Print()
}

// func TestMessagePrint(t *testing.T) {
// 	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
// 	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
// 	server := basePlatformSOMAS.CreateServer(m, 1, 1, time.Second)
// 	ag := NewTestAgent(server)
// 	newMsg := ag.NewTestMessage()
// 	newMsg.Print()
// 	originalStdout := os.Stdout
// 	r, w, _ := os.Pipe()
// 	os.Stdout = w
// 	newMsg.Print()
// 	w.Close()
// 	var buf bytes.Buffer
// 	io.Copy(&buf, r)
// 	output := buf.String()
// 	os.Stdout = originalStdout
// 	expected := "message received from " + ag.GetID().String() + "\n"
// 	if string(output) != string(expected) {
// 		t.Error("Expected", expected, "but got", output)
// 	}

// }

func (tba *TestServer) infMessageSend(newMsg InfiniteLoopMessage, receiver []uuid.UUID, done chan struct{}) {
	go tba.SendMessage(newMsg, receiver)
	tba.EndAgentListeningSession()
	done <- struct{}{}
}

func TestInfLoopProtection(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	timeLimit := 100 * time.Millisecond
	server := NewTestServer(m, 1, 1, timeLimit)
	ag1 := NewTestAgent(server)
	newMsg := CreateInfiniteLoopMessage()
	done := make(chan struct{}, 1)
	receiver := make([]uuid.UUID, 1)
	receiver[0] = ag1.GetID()
	go server.infMessageSend(newMsg, receiver, done)
	startTime := time.Now()
	time.Sleep(300 * time.Millisecond)
	select {
	case <-done:
		return
	default:
		timeTaken := time.Since(startTime)
		t.Error("Function did not terminate early on time limit. Time taken:", timeTaken, "expected:", timeLimit)
	}
}

type RunHandler struct {
	iters int
	turns int
}

func (r *RunHandler) RunRound() {
	r.iters += 1
}
func (r *RunHandler) RunTurn() {
	r.turns += 1
}

func TestGameRunner(t *testing.T) {
	m := make([]basePlatformSOMAS.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = basePlatformSOMAS.MakeAgentGeneratorCountPair(NewTestAgent, 1)
	timeLimit := 100 * time.Millisecond
	server := NewTestServer(m, 1, 1, timeLimit)
	runHandler := RunHandler{iters: 0, turns: 0}
	server.SetRunHandler(&runHandler)
	server.BaseServer.RunRound()
	server.BaseServer.RunTurn()
	if runHandler.iters != 1 {
		t.Errorf("Server unable to run round: have round value %d, expected %d", runHandler.iters, 1)
	}
	if runHandler.turns != 1 {
		t.Errorf("Server unable to run turn: have turn value %d, expected %d", runHandler.turns, 1)
	}
}
