package messaging

import (
	baseAgent "github.com/MattSScott/basePlatformSOMAS/pkg/agents/BaseAgent"
	message "github.com/MattSScott/basePlatformSOMAS/pkg/messaging"
	
)

func MessagingSession[T baseAgent.Agent](Agents []T) {
	var messageQueue []message.Message[T]
	messageCount:=0
	for _, agent := range Agents{
		messageQueue[messageCount]=message.CreateMessage(agent, "hello", Agents)
		messageCount ++
	}
	for _, message := range messageQueue {
		for _, recipient := range message.GetRecipients(){
			recipient.HandleMessage(message) 
			message.Messaging.HandleMessage(message)

		}
	}
}