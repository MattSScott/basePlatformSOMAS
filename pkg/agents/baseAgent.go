package agent

import (
	//"fmt"
	//"strconv"
	"github.com/google/uuid"
)

type Agent struct {
	ID uuid.UUID
	Name string
	Floor int
	Energy int
}

func NewAgent(name string, floor, energy int) Agent{
	return Agent{
		ID: uuid.New(),
		Name: name,
		Floor: floor,
		Energy: energy, //floor and energy in different type of agent
	}
}  

func ReduceEnergy(a *Agent){
	a.Energy-=10
}

func (a *Agent)UpdateAgent(){
		ReduceEnergy(a)
		
	}