package server_test

import (
	"sync"
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/internal/testUtils"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
	"github.com/google/uuid"
)

func TestGenerateServer(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != numAgents {
		t.Error(lenAgentMap, "agents initialised, expected:", numAgents)
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
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
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	server.Start()
}

func TestRunTurn(t *testing.T) {
	numAgents := 20
	iterations := 1
	turns := 1
	server := testUtils.GenerateTestServer(numAgents, iterations, turns, time.Millisecond, 100)
	server.SetGameRunner(server)
	server.Start()
	if server.TurnCounter != (iterations * turns) {
		t.Error("wrong number of iterations executed, got:", server.TurnCounter, "expected", iterations*turns)
	}
}

func TestDeliverMessage(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	for id, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
		server.DeliverMessage(testMessage, id)
	}
	server.ExposeEndListening()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestAddAgent(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	agent := testUtils.NewTestAgent(server)
	server.AddAgent(agent)
	agMap := server.GetAgentMap()
	lenAgMap := len(agMap)
	expectedNumAgents := numAgents + 1
	if lenAgMap != expectedNumAgents {
		t.Error("Removing Agents Failed,expected number of agents:", expectedNumAgents, ",got:", lenAgMap)
	}
}

func TestRemoveAgent(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	agMap := server.GetAgentMap()
	for _, ag := range agMap {
		server.RemoveAgent(ag)
	}
	lenAgMap := len(agMap)
	if lenAgMap != 0 {
		t.Error("Removing Agents Failed,expected number of agents: 0,got:", lenAgMap)
	}
}

func TestEndAgentListeningSession(t *testing.T) {
	numMessages := 20
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 5*time.Millisecond, 200)
	agentMap := server.GetAgentMap()
	agentGoal := int32(numMessages * numAgents)

	for _, ag := range agentMap {
		ag.SetGoal(agentGoal)
		msg := ag.CreateTestMessage()
		for recip := range agentMap {
			for i := 0; i < numMessages; i++ {
				ag.SendMessage(msg, recip)
			}
		}
	}

	start := time.Now()
	stat := server.ExposeEndListening()
	end := time.Since(start)
	if !stat {
		t.Error("Messaging ended on timeout, execution took:", end)
	}
	for _, ag := range agentMap {
		if !ag.ReceivedMessage() {
			t.Errorf("agent %s recieved %d messages, expected %d\n", ag.GetID(), ag.GetCounter(), ag.GetGoal())
		}
	}
}
func TestNumIterationsInServer(t *testing.T) {
	iterations := 1
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	getIterationsValue := server.GetIterations()
	if getIterationsValue != iterations {
		t.Error("Incorrect number of iterations instantiated, expected:", iterations, "got:", getIterationsValue)
	}
}

func TestNumTurnsInServer(t *testing.T) {
	turns := 1
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	getTurnsValue := server.GetTurns()
	if getTurnsValue != turns {
		t.Error("Incorrect number of turns instantiated, expected:", turns, "got:", getTurnsValue)
	}
}

func TestBroadcastMessage(t *testing.T) {
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100)
	agentGoal := int32(numAgents - 1)

	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(agentGoal)
		testMessage := ag.CreateTestMessage()
		ag.BroadcastMessage(testMessage)
	}
	server.ExposeEndListening()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestSynchronousMessagingSession(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	server.RunSynchronousMessagingSession()
	for _, ag := range server.GetAgentMap() {
		if ag.GetCounter() != int32(numAgents) {
			t.Error("All messages did not pass, got:", ag.GetCounter(), "expected:", numAgents)
		}
	}
}

func TestAccessAgentByID(t *testing.T) {
	var randNum int32 = 2357
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100)
	for _, ag := range server.GetAgentMap() {
		ag.SetCounter(randNum)
	}
	for id := range server.ViewAgentIdSet() {
		accessedAgentID := server.AccessAgentByID(id).GetCounter()
		if accessedAgentID != randNum {
			t.Error("Access Agent By ID is not working (incorrect struct value in test agent),expected:", randNum, "got:", accessedAgentID)
		}
	}
}

func TestGameRunner(t *testing.T) {
	timeLimit := 100 * time.Millisecond
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	server.SetGameRunner(server)
	server.RunTurn(-1, -1)
	turns := server.TurnCounter
	if turns != 1 {
		t.Errorf("Server unable to run turn: have turn value %d, expected %d", turns, 1)
	}
}

func TestIterationRunner(t *testing.T) {
	timeLimit := 100 * time.Millisecond
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	server.SetGameRunner(server)
	server.RunStartOfIteration(-1)
	server.RunEndOfIteration(-1)
	if server.IterationStartCounter != 1 {
		t.Errorf("RunStartOfIteration failed. Expected StartCounter to be 1, got: %d", server.IterationStartCounter)
	}
	if server.IterationEndCounter != 1 {
		t.Errorf("RunEndOfIteration failed. Expected EndCounter to be 1, got: %d", server.IterationEndCounter)
	}
}

func TestGoroutineWontHangAsyncMessaging(t *testing.T) {
	var counter uint32 = 0
	var numAgents int = 3
	timeLimit := 1000 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	wg := &sync.WaitGroup{}
	testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
	server.ExposeEndOfTurn()
	testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
	wg.Wait()
	goal := uint32(2 * numAgents)
	if counter != goal {
		t.Error(counter, "goroutines have exited,", goal, "were spawned")
	}
}

func TestRepeatedAsyncMessaging(t *testing.T) {
	var counter uint32 = 0
	var numAgents int = 3
	var numIters int = 5
	timeLimit := 100 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	for i := 0; i < numIters; i++ {
		wg := &sync.WaitGroup{}
		server.ExposeStartOfTurn()
		testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
		stat := server.ExposeEndListening()
		if !stat {
			t.Error("Session not correctly handled")
		}
		wg.Wait()
	}
	goal := uint32(numIters * numAgents)
	if counter != goal {
		t.Error(counter, "goroutines have exited,", goal, "were spawned")
	}
}

func TestTimeoutExit(t *testing.T) {
	var numAgents int = 3
	timeLimit := 100 * time.Millisecond
	agentWorkload := 150 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	server.ExposeStartOfTurn()
	timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
	timeoutMsg.SetSender(uuid.New())
	for _, ag := range server.GetAgentMap() {
		ag.BroadcastMessage(timeoutMsg)
	}
	status := server.ExposeEndListening()
	if status && (agentWorkload > timeLimit) {
		t.Error("Should have exited on timeout but did not")
	}
}

func TestRepeatedTimeouts(t *testing.T) {
	var numAgents int = 3
	var numIters int = 5
	timeLimit := 100 * time.Millisecond
	agentWorkload := 20 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
	timeoutMsg.SetSender(uuid.New())
	for i := 0; i < numIters; i++ {
		server.ExposeStartOfTurn()
		for _, ag := range server.GetAgentMap() {
			ag.BroadcastMessage(timeoutMsg)
		}
		status := server.ExposeEndListening()
		if status && (agentWorkload > timeLimit) {

			t.Error("Should have exited on timeout but did not", i)
		}
	}
}

func TestSendMessageNoIDPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("did not panic when message sender not set")
		}
	}()
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Millisecond, 100)
	agMap := server.GetAgentMap()

	for _, ag := range agMap {
		msg := &testUtils.TestMessage{}
		for recip := range agMap {
			ag.SendMessage(msg, recip)
		}
	}
}

func TestRunTurnNotSetPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("Did not panic when RunTurn not overriden in extended server")
		}
	}()
	server := &testUtils.TestTurnMethodPanics{
		BaseServer: server.CreateServer[testUtils.ITestBaseAgent](1, 1, time.Millisecond, 100),
	}
	server.RunTurn(0, 0)
}

func TestRunStartOfIterationNotSetPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("Did not panic when RunStartOfIteration not overriden in extended server")
		}
	}()
	server := &testUtils.TestTurnMethodPanics{
		BaseServer: server.CreateServer[testUtils.ITestBaseAgent](1, 1, time.Millisecond, 100),
	}
	server.RunStartOfIteration(0)
}

func TestRunEndOfIterationNotSetPanic(t *testing.T) {
	defer func() {
		if panicValue := recover(); panicValue == nil {
			t.Errorf("Did not panic when RunEndOfIteration not overriden in extended server")
		}
	}()
	server := &testUtils.TestTurnMethodPanics{
		BaseServer: server.CreateServer[testUtils.ITestBaseAgent](1, 1, time.Millisecond, 100),
	}
	server.RunEndOfIteration(0)
}

func TestRunStartEndOfTurn(t *testing.T) {
	timeLimit := 1 * time.Millisecond
	numAgents := 2
	iterations := 3
	server := testUtils.GenerateTestServer(numAgents, iterations, 1, timeLimit, 100)
	server.SetGameRunner(server)
	server.Start()

	if !(server.IterationStartCounter == iterations) {
		t.Error(server.IterationStartCounter, "instances of RunStartOfTurn(iteration) executed. Expected:", iterations)
	}
	if !(server.IterationEndCounter == iterations) {
		t.Error(server.IterationEndCounter, "instances of RunEndOfTurn(iteration) executed. Expected:", iterations)
	}
}

func TestMessagesSendInSaturatedServer(t *testing.T) {
	timeLimit := 1 * time.Millisecond
	server := testUtils.GenerateTestServer(0, 1, 1, timeLimit, 100)
	evilAgent1 := testUtils.NewTestAgent(server)
	evilAgent2 := testUtils.NewTestAgent(server)
	testAgent1 := testUtils.NewTestAgent(server)
	testAgent2 := testUtils.NewTestAgent(server)
	testAgent1.SetGoal(1)
	testAgent2.SetGoal(1)
	evilAgent1.SetGoal(0)
	evilAgent2.SetGoal(0)
	server.AddAgent(testAgent1)
	server.AddAgent(testAgent2)
	server.AddAgent(evilAgent1)
	server.AddAgent(evilAgent2)

	infLoopMessage := testUtils.CreateInfLoopMessage()
	infLoopMessage.SetSender(evilAgent1.GetID())
	evilAgent1.SendMessage(infLoopMessage, evilAgent2.GetID())
	time.Sleep(10 * time.Millisecond)
	//if message bandwidth is faulty this will fill it with messages from the two agents
	testMsg := testAgent1.CreateTestMessage()
	testAgent1.SendMessage(testMsg, testAgent2.GetID())
	testMsg = testAgent2.CreateTestMessage()
	testAgent2.SendMessage(testMsg, testAgent1.GetID())
	server.ExposeEndListening()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag.GetID(), "Got ", ag.GetCounter(), "messages", "expected:", ag.GetGoal())
		}

	}
}

func TestRecursiveInvokeMessageHandlerCalls(t *testing.T) {
	numAgents := 3
	timeLimit := time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100)
	msg := testUtils.CreateInfLoopMessage()
	for _, ag := range server.GetAgentMap() {
		msg.SetSender(ag.GetID())
		ag.BroadcastMessage(msg)
	}
	server.ExposeEndListening()
}

func TestSendMessage(t *testing.T) {
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	for id, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
		agent1.SendMessage(testMessage, id)
	}
	server.ExposeEndListening()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestBroadcastMessageFromAgent(t *testing.T) {
	numAgents := 3
	server := testUtils.GenerateTestServer(numAgents, 1, 1, 10*time.Millisecond, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	for _, ag := range server.GetAgentMap() {
		ag.SetGoal(1)
	}
	agent1.BroadcastMessage(testMessage)
	senderID := agent1.GetID()
	server.ExposeEndListening()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() && ag.GetID() != senderID {
			t.Error(ag, "Didn't Receive Message")
		} else if ag.ReceivedMessage() && ag.GetID() == senderID {
			t.Error(ag, "is sender and received its own message")
		}
	}
}
