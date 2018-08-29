package main

import (
	"flag"
	"fmt"
	"github.com/stevelacy/dhrt/pkg/core"
	"github.com/stevelacy/dhrt/pkg/server"
	"log"
	"time"
)

var version = "dev"
var flagDebug = flag.Bool("debug", false, "enable verbose debug mode")

func main() {
	flag.Parse()

	node := dhrt.Node{
		ListenAddr: "127.0.0.1",
		Port:       7453,
		Status:     "healthy",
	}
	nodes := []dhrt.Node{
		node,
	}

	config := dhrt.Config{
		Root:  "./data/", // if Root is 'memory' it will not save to disc
		Nodes: nodes,
		Node:  node,
	}
	// Open the database with config
	db, err := dhrt.Open(config)
	if err != nil {
		panic(err)
	}
	fmt.Printf("dhrt version %s\ninitializing\n", version)

	// Start the daemon
	if err := server.Start(db, config); err != nil {
		log.Fatal(err)
	}
	supervisor(db)
}

func supervisor(s *dhrt.Store) {
	timer := time.NewTicker(time.Second * 5) // every 5 seconds
	defer timer.Stop()
	for range timer.C {
		if *flagDebug {
			fmt.Printf("supervisor tick")
		}
		// Flush from memory to disc every 5 seconds
		// Check to see if other nodes are online every 30 seconds
	}
}
