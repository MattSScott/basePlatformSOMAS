package message

//baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"

type Message struct {
	sender     IAgentMessaging
	content    string
	recipients []IAgentMessaging
}

// create read-only message instance
func CreateMessage(sender IAgentMessaging, content string, recipients []IAgentMessaging) Message {
	return Message{
		sender:     sender,
		content:    content,
		recipients: recipients,
	}
}

func (m *Message) GetSender() IAgentMessaging {
	return m.sender
}

func (m *Message) GetContent() string {
	return m.content
}

func (m *Message) GetRecipients() []IAgentMessaging {
	return m.recipients
}
