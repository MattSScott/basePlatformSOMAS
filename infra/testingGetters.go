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

func (serv *BaseServer[T]) SetServerAgentChannel(id uuid.UUID,channelValue chan ServerNotification) {
	serv.serverAgentChannelMap[id] = channelValue
}
