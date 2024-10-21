package diagnosticsEngine_test

import (
	"testing"

	"github.com/MattSScott/basePlatformSOMAS/v2/internal/diagnosticsEngine"
)

func TestGetNumberSentMessages(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	initSends := engine.GetNumberSentMessages()
	if initSends != 0 {
		t.Error("Diagnostics engine intialised with non-zero number")
	}
	engine.ReportSendMessageStatus(true)
	sendsTrue := engine.GetNumberSentMessages()
	if sendsTrue != initSends+1 {
		t.Error("Diagnostics engine not correctly incremented with message")
	}
	engine.ReportSendMessageStatus(false)
	sendsFalse := engine.GetNumberSentMessages()
	if sendsFalse != sendsTrue+1 {
		t.Error("Diagnostics engine not correctly incremented with message")
	}
}

func TestGetNumberSuccesses(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	initSucc := engine.GetNumberMessageSuccesses()
	if initSucc != 0 {
		t.Error("Diagnostics engine intialised with non-zero number")
	}
	engine.ReportSendMessageStatus(true)
	newSuccTrue := engine.GetNumberMessageSuccesses()
	if newSuccTrue != initSucc+1 {
		t.Error("Diagnostics engine not incremented with success message")
	}
	engine.ReportSendMessageStatus(false)
	newSuccFalse := engine.GetNumberMessageSuccesses()
	if newSuccFalse != newSuccTrue {
		t.Error("Diagnostics engine incremented with failed message")
	}
}

func TestGetNumberEndMessagings(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	initEnds := engine.GetNumberEndMessagings()
	if initEnds != 0 {
		t.Error("Diagnostics engine intialised with non-zero number")
	}
	engine.ReportEndMessagingStatus(1)
	newEnds := engine.GetNumberEndMessagings()
	if newEnds != initEnds+1 {
		t.Error("Diagnostics engine successes not correctly incremented")
	}
}

func TestGetNumberDrops(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	numSends := 100
	numFails := 40
	for i := 0; i < numSends; i++ {
		if i < numFails {
			engine.ReportSendMessageStatus(false)
		} else {
			engine.ReportSendMessageStatus(true)
		}
	}
	trueDrops := engine.GetNumberMessageDrops()
	if trueDrops != numFails {
		t.Errorf("Incorrect number of dropped messages: have %d, expected %d\n", trueDrops, numFails)
	}
}

func TestMessagingSuccessRate(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	numSends := 100
	numSuccesses := 40
	for i := 0; i < numSends; i++ {
		if i < numSuccesses {
			engine.ReportSendMessageStatus(true)
		} else {
			engine.ReportSendMessageStatus(false)
		}
	}
	trueSuccessRate := engine.GetMessagingSuccessRate()
	expectedSuccessRate := float32(numSuccesses) / float32(numSends) * 100
	if trueSuccessRate != expectedSuccessRate {
		t.Errorf("Incorrect success rate: have %f, expected %f\n", trueSuccessRate, expectedSuccessRate)
	}
}

func TestEndMessagingSuccessRate(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	numAgents := 100
	numReports := 20

	engine.ReportEndMessagingStatus(numReports)
	trueSuccessRate := engine.GetEndMessagingSuccessRate(numAgents)
	expectedSuccessRate := float32(numReports) / float32(numAgents) * 100
	if trueSuccessRate != expectedSuccessRate {
		t.Errorf("Incorrect success rate: have %f, expected %f\n", trueSuccessRate, expectedSuccessRate)
	}
}

func TestResetDiagnostics(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	rounds := 3
	nMessages := 10
	for r := 0; r < rounds; r++ {
		for delta := 0; delta < nMessages; delta++ {
			engine.ReportSendMessageStatus(true)
		}
		engine.ReportEndMessagingStatus(nMessages)
		engine.ResetRoundDiagnostics()
		nSent := engine.GetNumberSentMessages()
		nSucc := engine.GetNumberMessageSuccesses()
		nEnds := engine.GetNumberEndMessagings()
		if nSent != 0 || nSucc != 0 || nEnds != 0 {
			t.Error("Diagnostic engine not reset at end of round")
		}
	}
}

func TestDivideByZeroProtection(t *testing.T) {
	engine := diagnosticsEngine.CreateDiagnosticsEngine()
	msgSuccessRate := engine.GetMessagingSuccessRate()
	endMsgingSuccessRate := engine.GetEndMessagingSuccessRate(0)
	if msgSuccessRate != 100.0 {
		t.Errorf("Diagnostic Engine incorrectly reported Message Success rate when 0 messages sent. Expected 100%%, got %v%%",msgSuccessRate)
	}
	if endMsgingSuccessRate != 100.0 {
		t.Errorf("Diagnostic Engine incorrectly reported Finished Messaging Success rate when 0 agents are present. Expected 100%%, got %v%%",endMsgingSuccessRate)
	}
}