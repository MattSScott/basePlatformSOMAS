package infra

import "github.com/google/uuid"

type PrivateServerFields[T IAgent[T]] interface {
	ServerNotification(uuid.UUID, ServerNotification)
	GetServerAgentChannel(uuid.UUID) chan ServerNotification
	SetServerAgentChannel(uuid.UUID, chan ServerNotification)
	GetAgentAgentChannel(uuid.UUID) chan IMessage[T]
	SetAgentAgentChannel(uuid.UUID, chan IMessage[T])
}

func (serv *BaseServer[T]) SendServerNotification(id uuid.UUID, serverNotification ServerNotification) {
	select {
	case serv.serverAgentChannelMap[id] <- serverNotification:
	default:
	}
}

func (serv *BaseServer[T]) GetServerAgentChannel(id uuid.UUID) chan ServerNotification {
	return serv.serverAgentChannelMap[id]
}

func (serv *BaseServer[T]) SetServerAgentChannel(id uuid.UUID, channelValue chan ServerNotification) {
	serv.serverAgentChannelMap[id] = channelValue
}

func (serv *BaseServer[T]) GetAgentAgentChannel(id uuid.UUID) chan IMessage[T] {
	return serv.agentAgentChannelMap[id]
}

func (serv *BaseServer[T]) SetAgentAgentChannel(id uuid.UUID, channelValue chan IMessage[T]) {
	serv.agentAgentChannelMap[id] = channelValue
}

func (serv *BaseServer[T]) BeginAgentListeningSession() {
	serv.beginAgentListeningSession()
}

// func (serv *BaseServer[T]) EndAgentListeningSession() {
// 	serv.endAgentListeningSession()
// }

func (serv *BaseServer[T]) Cleanup() {
	serv.cleanUp()
}

func (serv *BaseServer[T]) EndAgentListeningSession() {
	serv.endAgentListeningSession()
}

