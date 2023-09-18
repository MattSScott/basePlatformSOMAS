package message

//baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type Message struct {
	sender     Messaging
	content    string
	recipients []Messaging
}

// create read-only message instance
func CreateMessage(sender Messaging, content string, recipients []Messaging) Message {
	return Message{
		sender:     sender,
		content:    content,
		recipients: recipients,
	}
}

func (m *Message) GetSender() Messaging {
	return m.sender
}

func (m *Message) GetContent() string {
	return m.content
}

func (m *Message) GetRecipients() []Messaging {
	return m.recipients
}
