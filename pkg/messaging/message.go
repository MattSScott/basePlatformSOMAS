package message

import (
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
)

type Message struct {
	sender     baseAgent.Agent
	content    string
	recipients []baseAgent.Agent
}

func MessagingSession() {

}

// func CreateLetter() {

// }

// type LetterI[T baseAgent.Agent] interface { //better name would be great
// 	GetSnd() T
// 	SetSnd(a T)
// 	GetMsg() string
// 	SetMsg(s string)
// }

func CreateMessage(sender baseAgent.Agent, content string, recipients []baseAgent.Agent) Message {
	return Message{
		sender:     sender,
		content:    content,
		recipients: recipients,
	}
}

func (m *Message) GetSender() baseAgent.Agent {
	return m.sender
}

// func (l *Message[T]) SetSnd(a T) {
// 	l.sender = a
// }

func (m *Message) GetContent() string {
	return m.content
}

// func (l *Message[T]) SetMsg(s string) {
// 	l.message = s
// }

func (m *Message) GetRecipients() []baseAgent.Agent {
	return m.recipients
}

// func (l *Message[T]) SetRcv(a []T) {
// 	l.receivers = a
// }
