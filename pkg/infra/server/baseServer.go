package infra

import (
	"fmt"
	"strconv"
	"github.com/google/uuid"
)

type BaseServer struct {
	NumAgents int
	NumTurns  int
	Agents    []Agent
}

type Agent struct {
	ID uuid.UUID
	name string
	floor int
	energy int
}

func NewAgent(name string, floor, energy int) Agent{
	return Agent{
		ID: uuid.New(),
		name: name,
		floor: floor,
		energy: energy,
	}
}  

func ReduceEnergy(a *Agent){
	a.energy-=10
}

//function that updates the atributes of an agent
func (a *Agent)UpdateAgent(){
	ReduceEnergy(a)
	
}

func (bs *BaseServer) Init() {
	fmt.Println("Server Init")
	bs.NumAgents = 5
	bs.NumTurns = 4
	bs.Agents = make([]Agent, bs.NumAgents)
	for i := 0; i < bs.NumAgents; i++ {
		//converts the iteration to string
		name := strconv.Itoa(i)
		//creates new Agent
		bs.Agents[i]= NewAgent(name, 1, 100)
	}
}

func (bs *BaseServer) RunGameLoop(loopnum int) {
	fmt.Printf("Game Loop %d Running \n", loopnum)
	fmt.Printf("%d agents initialised: \n", bs.NumAgents)
	for index, element := range bs.Agents {
		//fmt.Printf("agent %d has id: %s \n", index, element.ID)
		fmt.Printf("agent %d has name: %s \n", index, element.name)
		//fmt.Printf("agent %d has name: %d \n", index, element.floor)
		fmt.Printf("agent %d has energy: %d \n", index, element.energy)
		fmt.Printf("_____________________________________________ \n")
		bs.Agents[index].UpdateAgent() 
	}


}

func (bs *BaseServer) Start() {
	bs.Init()
	//LOOPS
	for i:=0; i < bs.NumTurns; i++{
		fmt.Printf("Loop: %d \n", i)
		bs.RunGameLoop(i)
	}
	
}
