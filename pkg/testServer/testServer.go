package testserver

import (
	baseServer "somas_base_platform/pkg/infra/server"
)

type MyServer struct {
	*baseServer.BaseServer
	name      string
	numFloors int
}

func New(mapper []baseServer.AgentGeneratorCountPair) baseServer.Server {
	return MyServer{
		BaseServer: baseServer.CreateServer(mapper),
		name:       "Test Server",
		numFloors:  20,
	}
}
