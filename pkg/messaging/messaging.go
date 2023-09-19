package message

// base interface structure used for message passing - can be composed for more complex message structures
type IAgentMessaging interface {
	GetMessage() IMessage
	HandleMessage(IMessage)
}
