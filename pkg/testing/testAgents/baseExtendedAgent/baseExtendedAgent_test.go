package baseExtendedAgent_test

import (
	"testing"

	baseExtendedAgent "github.com/MattSScott/basePlatformSOMAS/pkg/testing/testAgents/baseExtendedAgent"
	"github.com/google/uuid"
	
	
)

func TestGetNewIExtendedAgent(t *testing.T) {
	iagent := baseExtendedAgent.GetNewIExtendedAgent("teststring")
	phrase := iagent.GetPhrase()
	id := iagent.GetID()

	if phrase != "teststring"{
		t.Error("Unsuccessful default implementation of IExtendedAgent due to string mismatch")
	}

	if id == uuid.Nil {
		t.Error("Unsuccessful default implementation of IExtendedAgent due to Nil ID")
	}
}