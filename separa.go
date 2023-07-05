package main

import (
	"separa/common/flag"
	"separa/common/log"
	"separa/core/run"
	"time"
)

func main() {
	startTime := time.Now()

	// First we parse the args
	flag.Parse()

	// Then we start the main process
	run.Start(&flag.Targets)

	// Finally we print the elapsed time
	elapsed := time.Since(startTime)
	log.Log.Printf("It costs %s", elapsed)
}
