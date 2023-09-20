package messaging

// // base interface structure used for message passing - can be composed for more complex message structures
type IAgentMessenger[T any] interface {
	// produces a list of messages (of any common interface) that an agent wishes to pass
	GetAllMessages([]T) []IMessage[T]
}
