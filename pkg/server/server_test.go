package server_test

import (
	"sync"
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/internal/testUtils"
	"github.com/google/uuid"
)

func TestGenerateServer(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != numAgents {
		t.Error(lenAgentMap, "agents initialised, expected:", numAgents)
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	server.Start()
}

func TestGenerateArrayFromMap(t *testing.T) {
	mapFound := make(map[uuid.UUID]int)
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
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
	numAgents := 20
	iterations := 1
	rounds := 1
	server := testUtils.GenerateTestServer(numAgents, iterations, rounds, time.Millisecond, 100000)
	server.SetGameRunner(server)
	server.Start()
	if server.GetTurnCounter() != (iterations * rounds) {
		t.Error("wrong number of iterations executed, got:", server.GetTurnCounter(), "expected", iterations*rounds)
	}
}

func TestRunIteration(t *testing.T) {
	numAgents := 2
	iterations := 2
	rounds := 2
	server := testUtils.GenerateTestServer(numAgents, iterations, rounds, time.Second, 100000)
	server.SetGameRunner(server)
	for i := 0; i < rounds; i++ {
		server.RunIteration()
	}
	if server.GetIterationCounter() != rounds {
		t.Error("wrong number of rounds executed", server.GetIterationCounter(), "expected", rounds)
	}
}

func TestAgentRecievesMessage(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	lenAgMap := len(server.GetAgentMap())
	arrayReceivers := make([]uuid.UUID, lenAgMap)
	i := 0
	for id, ag := range server.GetAgentMap() {
		arrayReceivers[i] = id
		i += 1
		ag.SetGoal(1)
	}
	server.SendMessage(testMessage, arrayReceivers)
	_ = server.EndAgentListeningSession()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestWaitForMessagingToEnd(t *testing.T) {
	numberOfMessages := 100
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Millisecond, 100000)
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id, ag := range agentMap {
		arrayOfIDs[i] = id
		i++
		ag.SetGoal(int32(numberOfMessages * numAgents))
	}

	for j := 0; j < numberOfMessages; j++ {
		for _, ag := range server.GetAgentMap() {
			msg := ag.CreateTestMessage()
			ag.SendMessage(msg, arrayOfIDs)
		}
	}
	start := time.Now()
	a := server.EndAgentListeningSession()
	end := time.Since(start)
	if !a {
		t.Error("Messaging ended on timeout,execution took:", end)
	}
	for _, ag := range agentMap {
		if !ag.ReceivedMessage() {
			t.Errorf("agent %s recieved %d messages, expected %d\n", ag.GetID(), ag.GetCounter(), ag.GetGoal())
		}
	}
}

func TestAddAgent(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
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
	iterations := 1
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	getIterationsValue := server.GetIterations()
	if getIterationsValue != iterations {
		t.Error("Incorrect number of iterations instantiated, expected:", iterations, "got:", getIterationsValue)
	}
}

func TestNumTurnsInServer(t *testing.T) {
	turns := 1
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	getTurnsValue := server.GetTurns()
	if getTurnsValue != turns {
		t.Error("Incorrect number of turns instantiated, expected:", turns, "got:", getTurnsValue)
	}
}

func TestBroadcastMessage(t *testing.T) {
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	i := 0
	var agentGoal int32 = int32(numAgents - 1)
	for _, ag := range server.GetAgentMap() {
		i += 1
		ag.SetGoal(agentGoal)
		testMessage := ag.CreateTestMessage()
		ag.BroadcastMessage(testMessage)
	}
	_ = server.EndAgentListeningSession()
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error(ag, "Didn't Receive Message")
		}
	}
}

func TestSendSynchronousMessage(t *testing.T) {
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	arrayReceivers := make([]uuid.UUID, numAgents)
	i := 0
	for id, ag := range server.GetAgentMap() {
		arrayReceivers[i] = id
		i += 1
		ag.SetGoal(1)
	}
	server.SendSynchronousMessage(testMessage, arrayReceivers)
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error("Didn't Receive Message")
		}
	}
}

func TestSynchronousMessagingSession(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
	server.RunSynchronousMessagingSession()
	for _, ag := range server.GetAgentMap() {

		if ag.GetCounter() != int32(numAgents-1) {
			t.Error("All messages did not pass, got:", ag.GetCounter(), "expected:", numAgents-1)
		}
	}
}

func TestAccessAgentByID(t *testing.T) {
	var randNum int32 = 2357
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second, 100000)
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

func TestMessagePrint(t *testing.T) {
	ag := testUtils.NewTestAgent(nil)
	msg := ag.CreateBaseMessage()
	msg.Print()
}

func TestGameRunner(t *testing.T) {
	timeLimit := 100 * time.Millisecond
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	server.SetGameRunner(server)
	server.BaseServer.RunIteration()
	server.BaseServer.RunTurn()
	turns := server.GetTurnCounter()
	iters := server.GetIterationCounter()
	if iters != 1 {
		t.Errorf("Server unable to run iteration: have round value %d, expected %d", iters, 1)
	}
	if turns != 1 {
		t.Errorf("Server unable to run turn: have turn value %d, expected %d", turns, 1)
	}
}

func TestGoroutineWontHangAsyncMessaging(t *testing.T) {
	var counter uint32 = 0
	var numAgents int = 3
	timeLimit := 1000 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	wg := &sync.WaitGroup{}
	testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
	server.HandleEndOfTurn(0, 0)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	for i := 0; i < numIters; i++ {
		wg := &sync.WaitGroup{}
		server.HandleStartOfTurn(0, i)
		testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
		stat := server.EndAgentListeningSession()
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	server.HandleStartOfTurn(0, 0)
	timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
	server.BroadcastMessage(timeoutMsg)
	status := server.EndAgentListeningSession()
	if status && (agentWorkload > timeLimit) {
		t.Error("Should have exited on timeout but did not")
	}
}

func TestRepeatedTimeouts(t *testing.T) {
	var numAgents int = 3
	var numIters int = 5
	timeLimit := 100 * time.Millisecond
	agentWorkload := 50 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	for i := 0; i < numIters; i++ {
		server.HandleStartOfTurn(0, i)
		timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
		server.BroadcastMessage(timeoutMsg)
		status := server.EndAgentListeningSession()
		if status && (agentWorkload > timeLimit) {
			t.Error("Should have exited on timeout but did not")
		}
	}
}

func TestRecursiveInvokeMessageHandlerCalls(t *testing.T) {
	numAgents := 3
	timeLimit := time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit, 100000)
	msg := testUtils.CreateInfLoopMessage()
	for _, ag := range server.GetAgentMap() {
		msg.SetSender(ag.GetID())
		ag.BroadcastMessage(msg)
	}
	server.EndAgentListeningSession()
}
