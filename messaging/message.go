package messaging

// base interface structure used for message - can be composed for more complex message structures
type IMessage[T any] interface {
	// returns the sender of a message
	GetSender() T
	// returns the list of agents that the message should be passed to
	GetRecipients() []T
	// calls the appropriate messsage handler method on the receiving agent
	InvokeMessageHandler(T)
}

// new message types can extend this
type BaseMessage[T IAgentMessenger[T]] struct {
	sender     T
	recipients []T
}

// create read-only message instance
func CreateMessage[T IAgentMessenger[T]](sender T, recipients []T) BaseMessage[T] {
	return BaseMessage[T]{
		sender:     sender,
		recipients: recipients,
	}
}

func (bm BaseMessage[T]) GetSender() T {
	return bm.sender
}

func (bm BaseMessage[T]) GetRecipients() []T {
	return bm.recipients
}

func (bm BaseMessage[T]) InvokeMessageHandler(agent T) {
}
