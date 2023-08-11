package baseagent
import (
	
	"github.com/google/uuid"
)

type Agent interface {
	Activity()
	UpdateAgent()
	GetID() uuid.UUID
	GetMsg() string 
	GetNet() []BaseAgent
	GetRcv() []BaseAgent
	SetMsg(s string) 
	SetNet(a []BaseAgent) 
	SetRcv(a []BaseAgent) 

}
