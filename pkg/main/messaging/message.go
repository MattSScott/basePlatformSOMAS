package messaging

// import baseagent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

// base interface structure used for message passing - can be composed for more complex message structures

// new message types extend this
type IMessage[T any] interface {
	GetSender() T
	GetRecipients() []T
	Accept(T)
}

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

// func CreateNullMessageWithSender(sender IAgentMessenger) BaseMessage {
// 	return BaseMessage{
// 		sender: sender,
// 	}
// }

func (bm BaseMessage[T]) GetSender() T {
	return bm.sender
}

func (bm BaseMessage[T]) GetRecipients() []T {
	return bm.recipients
}

func (bm BaseMessage[T]) Accept(agent T) {
}
