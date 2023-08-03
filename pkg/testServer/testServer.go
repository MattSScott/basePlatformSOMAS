package testserver

import (
	baseServer "basePlatformSOMAS/pkg/infra/server"
)

type MyServer struct {
	*baseServer.BaseServer
	name      string
	numFloors int //the user can change this
}

func New(mapper []baseServer.AgentGeneratorCountPair, numFloors int) baseServer.Server {
	return MyServer{
		BaseServer: baseServer.CreateServer(mapper),
		name:       "Test Server",
		numFloors:  numFloors,
	}
}
