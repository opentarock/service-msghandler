package main

import (
	"flag"
	"os"
	"os/signal"
	"runtime/pprof"

	"github.com/opentarock/service-api/go/client"
	"github.com/opentarock/service-api/go/log"
	nservice "github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-msghandler/routing"

	"github.com/opentarock/service-api/go/proto_msghandler"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	logger := log.New("name", "msghandler")
	flag.Parse()
	// profiliing related flag
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			logger.Error("Error creating cpuprofile file", "error", err)
			os.Exit(1)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	lobbyClient := client.NewLobbyClientNanomsg()
	err := lobbyClient.Connect("tcp://localhost:7001")
	if err != nil {
		logger.Error("Failed to connect to lobby service", "error", err)
	}
	defer lobbyClient.Close()

	msgHandlerService := nservice.NewRepService("tcp://*:11101")

	routeHandler := routing.NewRouteMessageHandler(lobbyClient)
	msgHandlerService.AddHandler(proto_msghandler.RouteMessageRequestType, routeHandler)

	err = msgHandlerService.Start()
	if err != nil {
		logger.Error("Error starting msghandler service", "error", err)
	}
	defer msgHandlerService.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	logger.Info("Interrupted", "signal", sig)
}
