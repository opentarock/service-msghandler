package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"

	nservice "github.com/opentarock/service-api/go/service"
	"github.com/opentarock/service-msghandler/routing"

	"github.com/opentarock/service-api/go/proto_msghandler"
)

var cpuprofile = flag.String("cpuprofile", "", "write cpu profile to file")

func main() {
	flag.Parse()
	// profiliing related flag
	if *cpuprofile != "" {
		f, err := os.Create(*cpuprofile)
		if err != nil {
			log.Fatal(err)
		}
		pprof.StartCPUProfile(f)
		defer pprof.StopCPUProfile()
	}

	msgHandlerService := nservice.NewRepService("tcp://*:11101")

	msgHandlerService.AddHandler(proto_msghandler.RouteMessageRequestType, routing.NewRouteMessageHandler())

	err := msgHandlerService.Start()
	if err != nil {
		log.Fatalf("Error starting message handler service: %s", err)
	}
	defer msgHandlerService.Close()

	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	sig := <-c
	log.Printf("Interrupted by %s", sig)
}
