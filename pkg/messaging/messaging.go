package message

type Messaging interface {
	GetMessage() Message
	HandleMessage()
}
