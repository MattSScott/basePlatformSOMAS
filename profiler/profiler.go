package profiler

import (
	"sync/atomic"
	"time"

	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/agent"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/message"
	"github.com/MattSScott/basePlatformSOMAS/v2/pkg/server"
)

type IProfilerAgent interface {
	agent.IAgent[IProfilerAgent]
	SetWorkload(int32)
	IncrementWorkload()
	FinishedWork() bool
	LogFinishedWork()
	SetStartTime()
	BroadcastSync(message.IMessage[IProfilerAgent])
	BroadcastAsync(message.IMessage[IProfilerAgent])
	CreateProfilerMessage() *ProfilerMessage
}

type IProfilerServer interface {
	server.IServer[IProfilerAgent]
}

type ProfilerAgent struct {
	*agent.BaseAgent[IProfilerAgent]
	workCompleted int32
	workload      int32
	startTime     time.Time
	endTime       time.Time
}

type ProfilerMessage struct {
	message.BaseMessage
}

func NewProfilerAgent(serv agent.IExposedServerFunctions[IProfilerAgent], agMessageGoal int32) IProfilerAgent {
	return &ProfilerAgent{
		BaseAgent: agent.CreateBaseAgent(serv),
		workload:  agMessageGoal,
	}
}

func (d *ProfilerMessage) InvokeMessageHandler(ag IProfilerAgent) {
	time.Sleep(10 * time.Millisecond)
	ag.IncrementWorkload()
	if ag.FinishedWork() {
		ag.LogFinishedWork()
	}
}

func (agent *ProfilerAgent) SetWorkload(workLoad int32) {
	agent.workload = workLoad
}

func (agent *ProfilerAgent) IncrementWorkload() {
	atomic.AddInt32(&agent.workCompleted, 1)
}

func (agent *ProfilerAgent) FinishedWork() bool {
	if agent.workCompleted == agent.workload {
		return true
	} else {
		return false
	}
}

func (agent *ProfilerAgent) LogFinishedWork() {
	agent.endTime = time.Now()
	agent.NotifyAgentFinishedMessaging()
}

func (agent *ProfilerAgent) SetStartTime() {
	agent.startTime = time.Now()
}

func (agent *ProfilerAgent) BroadcastSync(msg message.IMessage[IProfilerAgent]) {
	for id := range agent.ViewAgentIdSet() {
		if id == msg.GetSender() {
			continue
		}
		agent.SendSynchronousMessage(msg, id)
	}
}

func (agent *ProfilerAgent) BroadcastAsync(msg message.IMessage[IProfilerAgent]) {
	for id := range agent.ViewAgentIdSet() {
		if id == msg.GetSender() {
			continue
		}
		agent.SendMessage(msg, id)
	}
}

func (agent *ProfilerAgent) CreateProfilerMessage() *ProfilerMessage {
	return &ProfilerMessage{
		agent.CreateBaseMessage(),
	}
}
