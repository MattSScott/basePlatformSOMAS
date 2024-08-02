package infra

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BaseServer[U IExposedAgentFunctions, T IAgent[U]] struct {
	// map of agentid -> agent struct
	agentMap map[uuid.UUID]T
	// map of agentid -> empty struct so that agents cannot access each others agent structs
	agentIdSet map[uuid.UUID]struct{}
	// a map of agentid -> channel used by agents to send messages to agents
	agentAgentChannelMap map[uuid.UUID]chan IMessage
	// a map of agentid -> channel used by the server to send messages to agents
	serverAgentChannelMap map[uuid.UUID]chan ServerNotification
	//a channel that agents send their IDs to tell server they are ending messaging
	closureChannel chan uuid.UUID
	// waitground for server comms
	waitEnd *sync.WaitGroup
	// waitgroup for agent comms
	listeningWaitGroup *sync.WaitGroup
	// channel for server<->agent communication
	agentServerChannel chan uuid.UUID
	// map of uuid ->struct{}{} which stores the ids of agents which have stopped messaging
	agentStoppedTalkingMap map[uuid.UUID]struct{}
	// duration after which messaging phase forcefully ends during rounds
	maxMessagingDuration time.Duration
	// run single round
	roundRunner RoundRunner
	//iterations
	iterations int
}

func (server *BaseServer[T, U]) HandleStartOfTurn(iter, round int) {
	fmt.Printf("Iteration %d, Round %d starting...\n", iter, round)
	server.beginAgentListeningSession()

}

func (serv *BaseServer[T, U]) getAgentServerChannel() *chan uuid.UUID {
	return &serv.agentServerChannel
}

func (serv *BaseServer[T, U]) waitForMessagingToEnd() {
	//maxMessagingDuration := time.Second
	timeoutChannel := time.After(serv.maxMessagingDuration)

agentMessaging:
	for {
		select {

		case <-timeoutChannel:
			fmt.Println("len of stoppedtalkingmap:,", len(serv.agentStoppedTalkingMap))
			fmt.Println("Stopped messaging at time limit", serv.maxMessagingDuration, "seconds")
			break agentMessaging

		default:
			if len(serv.agentStoppedTalkingMap) == len(serv.agentMap) {
				fmt.Println("stopped messaging early", len(serv.agentStoppedTalkingMap), len(serv.agentMap))
				break agentMessaging
			}
			select {
			case uuid := <-serv.agentServerChannel:
				fmt.Println(uuid, "has stopped talking")
				serv.agentStoppedTalkingMap[uuid] = struct{}{}
			default:
				continue
			}
		}
	}
	serv.agentStoppedTalkingMap = make(map[uuid.UUID]struct{})
}

func (server *BaseServer[T, U]) HandleEndOfTurn(iter, round int) {
	server.waitForMessagingToEnd()
	server.endAgentListeningSession()
	fmt.Printf("Iteration %d, Round %d finished.\n", iter, round)
}

func (server *BaseServer[T, U]) RunAgentLoop() {

}

func (server *BaseServer[T, U]) SendMessage(msg IMessage, receiver uuid.UUID) {
	switch message := msg.(type) {
	case IMessage:
		select {
		case server.agentAgentChannelMap[receiver] <- message:
		default:
		}
	default:
		fmt.Println("unknown message type")
	}
}

func (serv *BaseServer[T, U]) AcknowledgeServerMessageReceived() {
	serv.listeningWaitGroup.Done()
}

func (server *BaseServer[T, U]) ReadChannel(agentID uuid.UUID) <-chan IMessage {
	return server.agentAgentChannelMap[agentID]
}

func (server *BaseServer[T, U]) initialiseChannels() {
	for _, agent := range server.agentMap {
		agentAgentChannel := make(chan IMessage, 20)
		serverAgentChannel := make(chan ServerNotification, 20)
		server.agentAgentChannelMap[agent.GetID()] = agentAgentChannel
		server.serverAgentChannelMap[agent.GetID()] = serverAgentChannel
	}
}

func (serv *BaseServer[U, T]) ViewAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[U, T]) AddAgent(agent T) {
	serv.agentMap[agent.GetID()] = agent
	serv.agentIdSet[agent.GetID()] = struct{}{}
}

func (serv *BaseServer[U, T]) ViewAgentIdSet() map[uuid.UUID]struct{} {
	return serv.agentIdSet
}

func (serv *BaseServer[U, T]) agentBeginSpin() {
	for _, agent := range serv.agentMap {
		serv.waitEnd.Add(1)
		agentAgentChannel := serv.agentAgentChannelMap[agent.GetID()]
		serverAgentChannel := serv.serverAgentChannelMap[agent.GetID()]
		go (agent).listenOnChannel(agentAgentChannel, serverAgentChannel, serv.waitEnd)
	}
}

func GenerateServer[U IExposedAgentFunctions, T IAgent[U]](maxDuration time.Duration, agentServerChannelBufferSize int) *BaseServer[U, T] {
	return &BaseServer[U, T]{
		agentMap:               make(map[uuid.UUID]T),
		agentIdSet:             make(map[uuid.UUID]struct{}),
		agentAgentChannelMap:   make(map[uuid.UUID]chan IMessage),
		serverAgentChannelMap:  make(map[uuid.UUID]chan ServerNotification),
		closureChannel:         make(chan uuid.UUID),
		waitEnd:                &sync.WaitGroup{},
		listeningWaitGroup:     &sync.WaitGroup{},
		agentStoppedTalkingMap: make(map[uuid.UUID]struct{}),
		agentServerChannel:     make(chan uuid.UUID, agentServerChannelBufferSize),
		maxMessagingDuration:   maxDuration,
		roundRunner:            nil, // TODO: need to initialise somehow (panic if uninitialised!)
	}
}

func (serv *BaseServer[T, U]) sendServerNotification(id uuid.UUID, serverNotification ServerNotification) {
	select {
	case serv.serverAgentChannelMap[id] <- serverNotification:
	default:
	}
}

func (serv *BaseServer[T, U]) beginAgentListeningSession() {
	fmt.Println("agents beginning to listen")
	for id := range serv.agentMap {
		serv.listeningWaitGroup.Add(1)
		serv.sendServerNotification(id, StartListeningNotification)
	}
	serv.listeningWaitGroup.Wait()
}

func (serv *BaseServer[T, U]) endAgentListeningSession() {
	fmt.Println("agents stopping listening")
	for id := range serv.agentMap {
		serv.listeningWaitGroup.Add(1)
		serv.sendServerNotification(id, EndListeningNotification)
	}
	serv.listeningWaitGroup.Wait()
}

func (serv *BaseServer[T, U]) Initialise() {
	serv.initialiseChannels()
	serv.agentBeginSpin()
}

func (serv *BaseServer[T, U]) Start() {
	serv.checkHandler()
	turns := 5
	iterations := 1
	for i := 0; i < iterations; i++ {
		for j := 0; j < turns; j++ {
			serv.HandleStartOfTurn(i+1, j+1)
			serv.roundRunner.RunTurn()
			serv.HandleEndOfTurn(i+1, j+1)
		}
	}
	serv.cleanUp()
}

func (serv *BaseServer[T, U]) cleanUp() {
	fmt.Println("Starting cleanup")
	for id := range serv.agentMap {
		serv.listeningWaitGroup.Add(1)
		select {
		case serv.serverAgentChannelMap[id] <- StopListeningSpinner:
		default:
		}
	}
	fmt.Println("Closure messages...")
	serv.awaitClosureMessages()
	serv.waitEnd.Wait()
	fmt.Println("closing channels")
	serv.closeChannels()
}

func (serv *BaseServer[T, U]) AcknowledgeClosure(id uuid.UUID) {
	//fmt.Println("sending id")
	serv.closureChannel <- id
}

func (serv *BaseServer[T, U]) closeChannels() {
	for _, channel := range serv.agentAgentChannelMap {
		close(channel)
	}
	for _, channel := range serv.serverAgentChannelMap {
		close(channel)
	}
	close(serv.closureChannel)
}

func (serv *BaseServer[T, U]) awaitClosureMessages() {
	fmt.Println("waiting for closures")
	closedAgents := make(map[uuid.UUID]struct{})
	for len(closedAgents) != len(serv.agentIdSet) {
		select {
		case id := <-serv.closureChannel:
			closedAgents[id] = struct{}{}
			fmt.Println("recieved")
		default:
			continue
		}
	}
}

func (serv *BaseServer[U, T]) GetAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[T, U]) agentStoppedTalking(id uuid.UUID) {
	//agentServerChannel := a.getAgentServerChannel()
	fmt.Println("sending stop talking request,id:", id)

	select {
	case serv.agentServerChannel <- id:
	default:
	}

}

func (serv *BaseServer[T, U]) SetRunHandler(handler RoundRunner) {
	serv.roundRunner = handler
}

func (serv *BaseServer[T, U]) checkHandler() {
	if serv.roundRunner == nil {
		panic("handler has not been set. Have you run SetRunHandler?")
	}
}

// func (bs *BaseServer[T, U]) GetAgentMap() map[uuid.UUID]T {
// 	return bs.agentMap
// }

// func (bs *BaseServer[T, U]) AddAgent(agentToAdd T) {
// 	bs.agentMap[agentToAdd.GetID()] = agentToAdd
// }

func (bs *BaseServer[T, U]) RemoveAgent(agentToAdd T) {
	delete(bs.agentMap, agentToAdd.GetID())
}

func (bs *BaseServer[T, U]) GetIterations() int {
	return bs.iterations
}

func (bs *BaseServer[T, U]) RunGameLoop() {
	for id, agent := range bs.agentMap {
		fmt.Printf("Agent %s updating state \n", id)
		agent.UpdateAgentInternalState()
	}
}

// func (bs *BaseServer[T, U]) Start() {
// 	fmt.Printf("Server initialised with %d agents \n", len(bs.agentMap))
// 	fmt.Print("\n")
// 	//LOOPS
// 	for i := 0; i < bs.iterations; i++ {
// 		fmt.Printf("Game Loop %d running... \n \n", i)
// 		fmt.Printf("Main game loop running... \n \n")
// 		bs.roundRunner.RunRound()
// 		fmt.Printf("\nMain game loop finished. \n \n")
// 		fmt.Printf("Messaging session started... \n \n")
// 		fmt.Printf("\nMessaging session completed \n \n")
// 		fmt.Printf("Game Loop %d completed. \n", i)
// 	}

// }

type AgentGenerator[U IExposedAgentFunctions, T IAgent[U]] func() T

type AgentGeneratorCountPair[U IExposedAgentFunctions, T IAgent[U]] struct {
	generator AgentGenerator[U, T]
	count     int
}

func MakeAgentGeneratorCountPair[U IExposedAgentFunctions, T IAgent[U]](generatorFunction AgentGenerator[U, T], count int) AgentGeneratorCountPair[U, T] {
	return AgentGeneratorCountPair[U, T]{
		generator: generatorFunction,
		count:     count,
	}
}

func (bs *BaseServer[U, T]) GenerateAgentArrayFromMap() []T {

	agentMapToArray := make([]T, len(bs.agentMap))

	i := 0
	for _, ag := range bs.agentMap {
		agentMapToArray[i] = ag
		i++
	}
	return agentMapToArray
}

func (bs *BaseServer[T, U]) SendSynchronousMessage(msg IMessage, recipients []uuid.UUID) {
	for _, recip := range recipients {
		if msg.GetSender() == recip {
			continue
		}
		msg.InvokeMessageHandler(recip)
	}

}

func (bs *BaseServer[T, U]) RunSynchronousMessagingSession() {
	for _, agent := range bs.agentMap {
		agent.RunSynchronousMessaging()
	}
}

// func (bs *BaseServer[T, U]) RunSynchronousMessagingSession() {
// 	for _, agent := range bs.agentMap {
// 		allMessages := agent.GetAllMessages()
// 		for _, msg := range allMessages {
// 			recipients := msg.GetRecipients()
// 			for _, recip := range recipients {
// 				if agent.GetID() == recip.GetID() {
// 					continue
// 				}
// 				msg.InvokeMessageHandler(recip)
// 			}
// 		}
// 	}
// }

func (bs *BaseServer[U, T]) initialiseAgents(m []AgentGeneratorCountPair[U, T]) {
	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agent := pair.generator()
			bs.AddAgent(agent)
		}
	}
}

func (bs *BaseServer[U, T]) CastAgentToExposedAgentFunctions(agent T) U {
	return agent.(IAgent[U])
}

func (bs *BaseServer[U, T]) GetAgentFromID(id uuid.UUID) U {
	fullAgent := bs.agentMap[id]
	return bs.CastAgentToExposedAgentFunctions(fullAgent)
}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[U IExposedAgentFunctions, T IAgent[U]](generatorArray []AgentGeneratorCountPair[U, T], iterations int) *BaseServer[U, T] {
	serv := &BaseServer[U, T]{
		agentMap:   make(map[uuid.UUID]T),
		iterations: iterations,
	}
	serv.initialiseAgents(generatorArray)
	return serv
}
