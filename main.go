package main

import (
	testserver "somas_base_platform/pkg/testServer"
)

func main() {
	// var s infra.Server = &infra.BaseServer{}
	// var s infra.Server = &testserver.MyServer{}
	// ts := testserver.MyServer{&infra.BaseServer{}}
	ts := testserver.New()
	// s.Start()
	ts.Init()
}
