package main

import (
	"time"

	"github.com/champly/lib4go/signal"
	"github.com/champly/mcpserver/server"
)

func main() {
	s := server.New(
		server.WithFreq(time.Second*1),
		server.WithAddress("0.0.0.0"),
		server.WithGRPCPort(20002),
	)

	s.Start(signal.SetupSignalHandler())
}
