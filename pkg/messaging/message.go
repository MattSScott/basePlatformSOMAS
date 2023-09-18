package message

import (
	//baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
)

type Message[T Agent] struct {
	sender     T
	content    string
	recipients []T
}
type Agent interface {
	GetMessage() Message[Agent]
	HandleMessage(Message[Agent]) Message[Agent]
	
}



func CreateMessage[T Agent](sender T, content string, recipients []T) Message[T] {
	return Message[T]{
		sender:     sender,
		content:    content,
		recipients: recipients,
	}
}
// func CreateMessageI[T baseAgent.Agent](sender T, content string, recipients []T) Messaging[T] {
// 	return &Message[T]{
// 		sender:     sender,
// 		content:    content,
// 		recipients: recipients,
// 	}
// }


func (m *Message[T]) GetSender() T {
	return m.sender
}



func (m *Message[T]) GetContent() string {
	return m.content
}



func (m *Message[T]) GetRecipients() []T {
	return m.recipients
}


