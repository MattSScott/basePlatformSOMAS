package messaging

// // base interface structure used for message passing - can be composed for more complex message structures
type IAgentMessenger[T any] interface {
	GetAllMessages([]T) []IMessage[T]
}
