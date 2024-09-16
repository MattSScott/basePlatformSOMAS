package server_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/internal/testUtils"
	"github.com/google/uuid"
)

func TestGenerateServer(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	lenAgentMap := len(server.GetAgentMap())
	if lenAgentMap != numAgents {
		t.Error(lenAgentMap, "agents initialised, expected:", numAgents)
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)

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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	server.Start()
}

func TestGenerateArrayFromMap(t *testing.T) {
	mapFound := make(map[uuid.UUID]int)
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
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
	server := testUtils.GenerateTestServer(numAgents, iterations, rounds, time.Millisecond)
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
	server := testUtils.GenerateTestServer(numAgents, iterations, rounds, time.Second)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
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
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
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
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	getIterationsValue := server.GetIterations()
	if getIterationsValue != iterations {
		t.Error("Incorrect number of iterations instantiated, expected:", iterations, "got:", getIterationsValue)
	}
}

func TestNumTurnsInServer(t *testing.T) {
	turns := 1
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	getTurnsValue := server.GetTurns()
	if getTurnsValue != turns {
		t.Error("Incorrect number of turns instantiated, expected:", turns, "got:", getTurnsValue)
	}
}

func TestSendSynchronousMessage(t *testing.T) {
	numAgents := 10
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	agent1 := testUtils.NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()

	arrayReceivers := make([]uuid.UUID, numAgents)
	i := 0
	for id, ag := range server.GetAgentMap() {
		arrayReceivers[i] = id
		i += 1
		ag.SetGoal(1)
	}

	server.EndAsyncMessaging()
	server.SendSynchronousMessage(testMessage, arrayReceivers)
	for _, ag := range server.GetAgentMap() {
		if !ag.ReceivedMessage() {
			t.Error("Didn't Receive Message")
		}
	}
}

func TestSynchronousMessagingSession(t *testing.T) {
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)
	server.EndAsyncMessaging()
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, time.Second)

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
	ag := testUtils.NewTestAgent(nil)
	msg := ag.CreateBaseMessage()
	msg.Print()
}

// func TestInfLoopProtection(t *testing.T) {
// 	server := testUtils.GenerateTestServer(1, 1, 1, 20*time.Millisecond)
// 	timeLimit := 100 * time.Millisecond

// 	ag1 := testUtils.NewTestAgent(server)
// 	newMsg := testUtils.CreateTestTimeoutMessage()
// 	done := make(chan struct{}, 1)
// 	receiver := make([]uuid.UUID, 1)
// 	receiver[0] = ag1.GetID()
// 	go server.InfMessageSend(newMsg, receiver, done)
// 	startTime := time.Now()
// 	select {
// 	case <-done:
// 		return
// 	case <-time.After(timeLimit):
// 		timeTaken := time.Since(startTime)
// 		t.Error("Function did not terminate early on time limit. Time taken:", timeTaken, "expected:", timeLimit)
// 	}
// }

func TestGameRunner(t *testing.T) {
	timeLimit := 100 * time.Millisecond
	numAgents := 2
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit)
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

	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit)
	wg := &sync.WaitGroup{}
	testUtils.SendNotifyMessages(server.GetAgentMap(), &counter, wg)
	server.EndAgentListeningSession()
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

	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit)
	numAgentsInt32 := uint32(numAgents)
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

	goal := uint32(numIters) * numAgentsInt32
	if counter != goal {
		t.Error(counter, "goroutines have exited,", goal, "were spawned")
	}
}

func TestTimeoutExit(t *testing.T) {
	var numAgents int = 3
	timeLimit := 100 * time.Millisecond

	agentWorkload := 150 * time.Millisecond
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit)
	server.HandleStartOfTurn(0,0)
	timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
	server.BroadcastMessage(&timeoutMsg)
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
	server := testUtils.GenerateTestServer(numAgents, 1, 1, timeLimit)
	for i := 0; i < numIters; i++ {
		server.HandleStartOfTurn(0,i)
		timeoutMsg := testUtils.CreateTestTimeoutMessage(agentWorkload)
		server.BroadcastMessage(&timeoutMsg)
		status := server.EndAgentListeningSession()
		if status && (agentWorkload > timeLimit) {
			t.Error("Should have exited on timeout but did not")
		}
	}
}