package infra

import (
	"fmt"
	"sync"
	"time"

	"github.com/google/uuid"
)

type BaseServer[T IAgent[T]] struct {
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

func (server *BaseServer[T]) HandleStartOfTurn(iter, round int) {
	fmt.Printf("Iteration %d, Round %d starting...\n", iter, round)
	server.beginAgentListeningSession()

}

func (serv *BaseServer[T]) waitForMessagingToEnd() {
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

func (server *BaseServer[T]) HandleEndOfTurn(iter, round int) {
	server.waitForMessagingToEnd()
	server.endAgentListeningSession()
	fmt.Printf("Iteration %d, Round %d finished.\n", iter, round)
}

func (server *BaseServer[T]) RunAgentLoop() {

}

func (server *BaseServer[T]) SendMessage(msg IMessage, receivers []uuid.UUID) {
	for _,receiver := range receivers {
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
}

func (serv *BaseServer[T]) AcknowledgeServerMessageReceived() {
	serv.listeningWaitGroup.Done()
}

func (server *BaseServer[T]) ReadChannel(agentID uuid.UUID) <-chan IMessage {
	return server.agentAgentChannelMap[agentID]
}

func (server *BaseServer[T]) initialiseChannels() {
	for _, agent := range server.agentMap {
		agentAgentChannel := make(chan IMessage, 20)
		serverAgentChannel := make(chan ServerNotification, 20)
		server.agentAgentChannelMap[agent.GetID()] = agentAgentChannel
		server.serverAgentChannelMap[agent.GetID()] = serverAgentChannel
	}
}

func (serv *BaseServer[T]) ViewAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

func (serv *BaseServer[T]) AddAgent(agent T) {
	serv.agentMap[agent.GetID()] = agent
	serv.agentIdSet[agent.GetID()] = struct{}{}
}

func (serv *BaseServer[T]) ViewAgentIdSet() map[uuid.UUID]struct{} {
	return serv.agentIdSet
}

func (serv *BaseServer[T]) agentBeginSpin() {
	for _, agent := range serv.agentMap {
		serv.waitEnd.Add(1)
		agentAgentChannel := serv.agentAgentChannelMap[agent.GetID()]
		serverAgentChannel := serv.serverAgentChannelMap[agent.GetID()]
		go agent.listenOnChannel(agentAgentChannel, serverAgentChannel, serv.waitEnd)
	}
}

func (serv *BaseServer[T]) AccessAgentByID(id uuid.UUID) T {
	return serv.agentMap[id]
}

func GenerateServer[T IAgent[T]](maxDuration time.Duration, agentServerChannelBufferSize int) *BaseServer[T] {
	return &BaseServer[T]{
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

func (serv *BaseServer[T]) sendServerNotification(id uuid.UUID, serverNotification ServerNotification) {
	select {
	case serv.serverAgentChannelMap[id] <- serverNotification:
	default:
	}
}

func (serv *BaseServer[T]) beginAgentListeningSession() {
	fmt.Println("agents beginning to listen")
	for id := range serv.agentMap {
		serv.listeningWaitGroup.Add(1)
		serv.sendServerNotification(id, StartListeningNotification)
	}
	serv.listeningWaitGroup.Wait()
}

func (serv *BaseServer[T]) endAgentListeningSession() {
	fmt.Println("agents stopping listening")
	for id := range serv.agentMap {
		serv.listeningWaitGroup.Add(1)
		serv.sendServerNotification(id, EndListeningNotification)
	}
	serv.listeningWaitGroup.Wait()
}

func (serv *BaseServer[T]) Initialise() {
	serv.initialiseChannels()
	serv.agentBeginSpin()
}

func (serv *BaseServer[T]) Start() {
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

func (serv *BaseServer[T]) cleanUp() {
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

func (serv *BaseServer[T]) AcknowledgeClosure(id uuid.UUID) {
	//fmt.Println("sending id")
	serv.closureChannel <- id
}

func (serv *BaseServer[T]) closeChannels() {
	for _, channel := range serv.agentAgentChannelMap {
		close(channel)
	}
	for _, channel := range serv.serverAgentChannelMap {
		close(channel)
	}
	close(serv.closureChannel)
}

func (serv *BaseServer[T]) awaitClosureMessages() {
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

func (serv *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
	return serv.agentMap
}

// TODO: this is never called
func (serv *BaseServer[T]) agentStoppedTalking(id uuid.UUID) {
	fmt.Println("sending stop talking request,id:", id)
	select {
	case serv.agentServerChannel <- id:
	default:
	}
}

func (serv *BaseServer[T]) SetRunHandler(handler RoundRunner) {
	serv.roundRunner = handler
}

func (serv *BaseServer[T]) checkHandler() {
	if serv.roundRunner == nil {
		panic("handler has not been set. Have you run SetRunHandler?")
	}
}

func (serv *BaseServer[T]) RunTurn() {}

// func (bs *BaseServer[T]) GetAgentMap() map[uuid.UUID]T {
// 	return bs.agentMap
// }

// func (bs *BaseServer[T]) AddAgent(agentToAdd T) {
// 	bs.agentMap[agentToAdd.GetID()] = agentToAdd
// }

func (bs *BaseServer[T]) RemoveAgent(agentToAdd T) {
	delete(bs.agentMap, agentToAdd.GetID())
}

func (bs *BaseServer[T]) GetIterations() int {
	return bs.iterations
}

func (bs *BaseServer[T]) RunGameLoop() {
	if bs.roundRunner == nil {
		panic("roundRuner has not been set.")

	}
	for id, agent := range bs.agentMap {
		fmt.Printf("Agent %s updating state \n", id)
		agent.UpdateAgentInternalState()
	}
}

func (bs *BaseServer[T]) RunRound() {}

// func (bs *BaseServer[T]) Start() {
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

type AgentGenerator[T IAgent[T]] func() T

type AgentGeneratorCountPair[T IAgent[T]] struct {
	generator AgentGenerator[T]
	count     int
}

func MakeAgentGeneratorCountPair[T IAgent[T]](generatorFunction AgentGenerator[T], count int) AgentGeneratorCountPair[T] {
	return AgentGeneratorCountPair[T]{
		generator: generatorFunction,
		count:     count,
	}
}

func (bs *BaseServer[T]) GenerateAgentArrayFromMap() []T {

	agentMapToArray := make([]T, len(bs.agentMap))

	i := 0
	for _, ag := range bs.agentMap {
		agentMapToArray[i] = ag
		i++
	}
	return agentMapToArray
}

func (bs *BaseServer[T]) SendSynchronousMessage(msg IMessage, recipients []uuid.UUID) {
	for _, recip := range recipients {
		if msg.GetSender() == recip {
			continue
		}
		msg.InvokeMessageHandler(recip)
	}

}

func (bs *BaseServer[T]) RunSynchronousMessagingSession() {
	for _, agent := range bs.agentMap {
		agent.RunSynchronousMessaging()
	}
}

func (bs *BaseServer[T]) initialiseAgents(m []AgentGeneratorCountPair[T]) {

	for _, pair := range m {
		for i := 0; i < pair.count; i++ {
			agent := pair.generator()
			bs.AddAgent(agent)
		}
	}

}

// generate a server instance based on a mapping function and number of iterations
func CreateServer[T IAgent[T]](generatorArray []AgentGeneratorCountPair[T], iterations int) *BaseServer[T] {
	serv := &BaseServer[T]{
		agentMap:   make(map[uuid.UUID]T),
		iterations: iterations,
	}
	serv.initialiseAgents(generatorArray)
	return serv
}
