package message

import (
	baseAgent "basePlatformSOMAS/pkg/agents/BaseAgent"
)

type Message[T baseAgent.Agent] struct {
	sender     T
	content    string
	recipients []T
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

func CreateMessage[T baseAgent.Agent](sender T, content string, recipients []T) Message[T] {
	return Message[T]{
		sender:     sender,
		content:    content,
		recipients: recipients,
	}
}

func (m *Message[T]) GetSender() T {
	return m.sender
}

// func (l *Message[T]) SetSnd(a T) {
// 	l.sender = a
// }

func (m *Message[T]) GetContent() string {
	return m.content
}

// func (l *Message[T]) SetMsg(s string) {
// 	l.message = s
// }

func (m *Message[T]) GetRecipients() []T {
	return m.recipients
}

// func (l *Message[T]) SetRcv(a []T) {
// 	l.receivers = a
// }
