package infra

import "github.com/google/uuid"

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

func (serv *BaseServer[T]) GetAgentAgentChannel(id uuid.UUID) chan IMessage {
	return serv.agentAgentChannelMap[id]
}

func (serv *BaseServer[T]) SetAgentAgentChannel(id uuid.UUID, channelValue chan IMessage) {
	serv.agentAgentChannelMap[id] = channelValue
}
