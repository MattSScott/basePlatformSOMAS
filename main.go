package main

import (
	testserver "somas_base_platform/pkg/testServer"
)

func main() {
	ts := testserver.New()
	ts.Init()
	ts.Start()
}
