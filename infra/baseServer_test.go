package infra_test

import (
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/infra"
	"github.com/google/uuid"
)

type ITestBaseAgent interface {
	infra.IAgent[ITestBaseAgent]
	CreateTestMessage() TestMessage
	HandleTestMessage()
	ReceivedMessage() bool
}

type ITestServer interface {
	infra.IServer[ITestBaseAgent]
	infra.PrivateServerFields[ITestBaseAgent]
}

type TestAgent struct {
	*infra.BaseAgent[ITestBaseAgent]
	receivedMessage bool
}

type TestAgent2 struct {
	*infra.BaseAgent[ITestBaseAgent]
	counter int
	limit   int
	goal    int
}

type TestServer struct {
	*infra.BaseServer[ITestBaseAgent]
}

type TestMessage struct {
	infra.BaseMessage
	arbitraryField int
}

func (tm TestMessage) InvokeMessageHandler(ag ITestBaseAgent) {
	ag.HandleTestMessage()
}

func (tba TestAgent) CreateTestMessage() TestMessage {
	return TestMessage{
		infra.BaseMessage{},
		5,
	}
}

func (tba TestAgent2) CreateTestMessage() TestMessage {
	return TestMessage{
		infra.BaseMessage{},
		5,
	}
}

func NewTestMessage() TestMessage {
	return TestMessage{
		infra.BaseMessage{},
		5,
	}
}

func NewTestAgent(serv infra.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestAgent{
		BaseAgent:       infra.CreateBaseAgent(serv),
		receivedMessage: false,
	}
}

func NewTestAgent2(serv infra.IExposedServerFunctions[ITestBaseAgent]) ITestBaseAgent {
	return &TestAgent2{
		BaseAgent: infra.CreateBaseAgent(serv),
		counter:   0,
		goal:      0,
	}
}

func (ag *TestAgent) HandleTestMessage() {
	ag.receivedMessage = true
	ag.NotifyAgentFinishedMessaging()
	ag.SetListeningSpinner(true)
}

func (a *TestAgent2) setLimit(limit int) {
	a.limit = limit
}

func (ag *TestAgent2) HandleTestMessage() {
	ag.counter += 1
	if ag.counter == ag.limit {
		ag.NotifyAgentFinishedMessaging()
		ag.SetListeningSpinner(true)
	}

}

func (ag *TestAgent2) ReceivedMessage() bool {
	if ag.counter == ag.goal {
		return true
	} else {
		return false
	}

}

func (ag *TestAgent) ReceivedMessage() bool {
	return ag.receivedMessage
}

// func NewTestBaseAgent() ITestBaseAgent {
// 	serv := &infra.BaseServer[ITestBaseAgent]{}

// 	return &TestAgent{
// 		BaseAgent: infra.CreateBaseAgent[ITestBaseAgent](serv),
// 		receivedMessage: false,
// 	}
// }

// func NewTestAgent1(serv *infra.BaseServer[ITestBaseAgent]) *TestAgent {
// 	serv = infra.GenerateServer[ITestBaseAgent](time.Second,2)
// 	return &TestAgent{
// 		BaseAgent:       infra.CreateBaseAgent[ITestBaseAgent](serv),
// 		receivedMessage: false,
// 	}
// }

func TestGenerateServer(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := infra.CreateServer[ITestBaseAgent](m, 1, time.Second, 1000, 20)
	if len(server.GetAgentMap()) != 3 {
		t.Error("len of agentmap is ", len(server.GetAgentMap()))
	}
}

func TestAgentsCorrectlyInstantiated(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, 3)

	server := infra.CreateServer(m, 1, time.Second, 1000, 20)

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
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := infra.CreateServer(m, 1, time.Second, 1000, 20)
	server.Initialise()
	server.RunGameLoop()
}

func TestSpinStart(t *testing.T) {
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := infra.CreateServer(m, 1, time.Second, 1000, 20)
	server.Initialise()
	arbitraryAgentID := uuid.New()
	server.SetServerAgentChannel(arbitraryAgentID, make(chan infra.ServerNotification, 1))
	//create a fake entry in the serverAgentChannelMap to send messages
	// to that wont be checked by an agent
	var waitGroup sync.WaitGroup
	waitGroup.Add(1)

	go func() {
		defer waitGroup.Done()
		server.SendServerNotification(arbitraryAgentID, infra.StartListeningNotification)

	}()

	waitGroup.Wait()
	msg := <-server.GetServerAgentChannel(arbitraryAgentID)
	if msg != infra.StartListeningNotification {
		t.Errorf("Incorrect Message Sent")
	}

}

func TestAgentAgentMessage(t *testing.T) {
	//server := infra.GenerateServer[ITestBaseAgent](time.Second, 2)
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, 3)
	server := infra.CreateServer(m, 1, time.Second, 1000, 20)

	arbitraryAgentID := uuid.New()

	server.SetAgentAgentChannel(arbitraryAgentID, make(chan infra.IMessage[ITestBaseAgent]))

	agent1 := NewTestAgent(server)
	testMessage := agent1.CreateTestMessage()
	server.AddAgent(agent1)
	arrayReceivers := make([]uuid.UUID, 1)
	arrayReceivers[0] = arbitraryAgentID
	//server.SendMessage(testMessage,arrayReceivers)
	server.Initialise()
	go func() {

		//defer waitGroup.Done()
		agent1.SendMessage(testMessage, arrayReceivers)
	}()

	//waitGroup.Wait()
	msg := <-server.GetAgentAgentChannel(arbitraryAgentID)
	if _, ok := msg.(TestMessage); ok {
		t.Logf("Message sent")
	} else {
		t.Errorf("message not sent")
	}
}

func TestAgentListeningSpinner(t *testing.T) {
	numberOfMessages := 10000
	const numAgents = 10
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent2, numAgents)
	server := infra.CreateServer(m, 1, 5*time.Second, numberOfMessages, 20)
	server.Initialise()

	//msg := NewTestMessage()
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id, ag := range agentMap {
		arrayOfIDs[i] = id
		i++
		agent := ag.(*TestAgent2)
		agent.setLimit(numberOfMessages)
		agent.goal = numberOfMessages
	}
	server.BeginAgentListeningSession()
	fmt.Println("ENDED BEGIN AGENT LISTENING")
	//j = j + 1

	for j := 0; j < numberOfMessages; j++ {
		msg := NewTestMessage()
		server.SendMessage(msg, arrayOfIDs)
	}
	//time.Sleep(5*time.Second)
	fmt.Println("FINISHED SENDING MESSAGES")
	server.EndAgentListeningSession()
	//server.WaitForMessagingToEnd()
	//server.WaitForMessagingToEnd()
	for _, ag := range agentMap {
		agent := ag.(*TestAgent2)
		if !agent.ReceivedMessage() {
			t.Errorf("agent %s recieved %d messages, expected %d\n", agent.GetID(), agent.counter, agent.goal)
		}
	}
	//t.Errorf("ended function")
}

func TestAgentListeningSpinnerOpen(t *testing.T) {

	const numAgents = 10
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := infra.CreateServer(m, 1, 5*time.Second, 1000, 20)
	server.Initialise()
	msg := NewTestMessage()
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id := range agentMap {
		arrayOfIDs[i] = id
		i++
	}
	server.BeginAgentListeningSession()
	fmt.Println("ENDED BEGIN AGENT LISTENING")
	//j = j + 1
	server.SendMessage(msg, arrayOfIDs)
	//time.Sleep(5*time.Second)
	fmt.Println("FINISHED SENDING MESSAGES")
	server.EndAgentListeningSession()
	//server.WaitForMessagingToEnd()
	//server.WaitForMessagingToEnd()
	for _, ag := range agentMap {
		agent := ag.(*TestAgent)
		if !agent.ReceivedMessage() {
			t.Errorf("agent didnt recieve a message")
		}
	}
	//t.Errorf("ended function")
}

func TestAgentListeningSpinnerClosed(t *testing.T) {
	const numAgents = 3
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := infra.CreateServer(m, 1, time.Second, 1000, 20)
	server.Initialise()
	msg := NewTestMessage()
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id := range agentMap {
		arrayOfIDs[i] = id
		i++
	}
	server.BeginAgentListeningSession()
	for _,ag := range agentMap {
		ag.SetListeningSpinner(true)
	}
	server.SendMessage(msg, arrayOfIDs)
	fmt.Println(arrayOfIDs)
	//time.Sleep(10*time.Second)

	for _, ag := range agentMap {
		testAgent := ag.(*TestAgent)
		if testAgent.receivedMessage != false {
			t.Errorf("agent did not stop listening")
		}
	}
}

func TestAgentChannelsClosed(t *testing.T) {

	const numAgents = 3
	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
	m[0] = infra.MakeAgentGeneratorCountPair(NewTestAgent, numAgents)
	server := infra.CreateServer(m, 1, time.Second, 1000, 20)
	server.Initialise()
	agentMap := server.GetAgentMap()
	arrayOfIDs := make([]uuid.UUID, numAgents)
	i := 0
	for id := range agentMap {
		arrayOfIDs[i] = id
		i++
	}
	server.BeginAgentListeningSession()
	server.EndAgentListeningSession()
	server.Cleanup()
	for id := range agentMap {
		agentChannel := server.GetAgentAgentChannel(id)
		_,ok := <-agentChannel
		if ok {
			t.Errorf("Channel not closed %s",id)
		}
	}
}	
	// fmt.Println(arrayOfIDs)
	// //time.Sleep(10*time.Second)

	// for _, ag := range agentMap {
	// 	testAgent := ag.(*TestAgent)
	// 	if testAgent.receivedMessage != false {
	// 		t.Errorf("agent did not stop listening")
	// 	}



// func TestNumIterationsInServer(t *testing.T) {
// 	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
// 	m[0] = infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

// 	server := infra.CreateServer[ITestBaseAgent](m, 1)

// 	if server.GetIterations() != 1 {
// 		t.Error("Incorrect number of iterations instantiated")
// 	}

// }

// type IExtendedTestServer interface {
// 	infra.IServer[ITestBaseAgent]
// 	GetAdditionalField() int
// }

// type ExtendedTestServer struct {
// 	*infra.BaseServer[ITestBaseAgent]
// 	testField int
// }

// func (ets *ExtendedTestServer) GetAdditionalField() int {
// 	return ets.testField
// }

// func (ets *ExtendedTestServer) RunGameLoop() {

// 	ets.BaseServer.RunGameLoop()
// 	ets.testField += 1
// }

// func CreateTestServer(mapper []infra.AgentGeneratorCountPair[ITestBaseAgent], iters int) IExtendedTestServer {
// 	return &ExtendedTestServer{
// 		BaseServer: infra.CreateServer[ITestBaseAgent](mapper, iters),
// 		testField:  0,
// 	}
// }

// func TestAddAgent(t *testing.T) {

// 	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

// 	agent1 := infra.NewBaseAgent[ITestBaseAgent]()

// 	baseServer.AddAgent(agent1)

// 	if len(baseServer.GetAgentMap()) != 1 {
// 		t.Error("Agent not correctly added to map")
// 	}
// }

// func TestRemoveAgent(t *testing.T) {

// 	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)

// 	agent1 := infra.NewBaseAgent[ITestBaseAgent]()

// 	baseServer.AddAgent(agent1)
// 	baseServer.RemoveAgent(agent1)

// 	if len(baseServer.GetAgentMap()) != 0 {
// 		t.Error("Agent not correctly removed from map")
// 	}
// }

// func TestFullAgentHashmap(t *testing.T) {
// 	baseServer := infra.CreateServer[ITestBaseAgent]([]infra.AgentGeneratorCountPair[ITestBaseAgent]{}, 1)
// 	for i := 0; i < 5; i++ {
// 		baseServer.AddAgent(infra.NewBaseAgent[ITestBaseAgent]())
// 	}

// 	for id, agent := range baseServer.GetAgentMap() {
// 		if agent.GetID() != id {
// 			t.Error("Server agent hashmap key doesn't match object")
// 		}
// 	}
// }

// func TestServerGameLoop(t *testing.T) {
// 	m := make([]infra.AgentGeneratorCountPair[ITestBaseAgent], 1)
// 	m[0] = infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

// 	server := CreateTestServer(m, 1)

// 	if server.GetAdditionalField() != 0 {
// 		t.Error("Additional server parameter not correctly instantiated")
// 	}

// 	server.RunGameLoop()

// 	if server.GetAdditionalField() != 1 {
// 		t.Error("Run Game Loop method not successfully overridden")
// 	}

// }

// func TestServerStartsCorrectly(t *testing.T) {
// 	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

// 	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

// 	baseServer.Start()
// }

// func TestAgentMapConvertsToArray(t *testing.T) {
// 	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

// 	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

// 	if len(baseServer.GenerateAgentArrayFromMap()) != 3 {
// 		t.Error("Agents not correctly mapped to array")
// 	}
// }

// func (tba *TestBaseAgent) GetAllMessages(others []ITestBaseAgent) []infra.IMessage[ITestBaseAgent] {
// 	msg := infra.CreateMessage[ITestBaseAgent](tba, others)

// 	return []infra.IMessage[ITestBaseAgent]{msg}
// }

// func TestMessagingSession(t *testing.T) {
// 	generator := infra.MakeAgentGeneratorCountPair[ITestBaseAgent](NewTestBaseAgent, 3)

// 	baseServer := infra.CreateServer([]infra.AgentGeneratorCountPair[ITestBaseAgent]{generator}, 1)

// 	agentArray := baseServer.GenerateAgentArrayFromMap()

// 	agent1 := agentArray[0]

// 	messages := agent1.GetAllMessages(agentArray)

// 	for _, msg := range messages {
// 		if len(msg.GetRecipients()) != 3 {
// 			t.Error("Incorrect number of message recipients")
// 		}
// 		for _, recip := range msg.GetRecipients() {
// 			if recip.GetID() == agent1.GetID() {
// 				continue
// 			}
// 			msg.InvokeMessageHandler(recip)
// 		}
// 	}

// }
