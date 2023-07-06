package main

import (
	"embed"
	"github.com/lcvvvv/appfinger"
	"separa/common/flag"
	"separa/common/log"
	"separa/core/run"
	"time"
)

//go:embed static/fingerprint.txt
var fingerprintEmbed embed.FS

const fingerprintPath = "static/fingerprint.txt"

func main() {
	startTime := time.Now()
	// initialize the fingerprint
	fs, _ := fingerprintEmbed.Open(fingerprintPath)
	if n, err := appfinger.InitDatabaseFS(fs); err != nil {
		log.Log.Fatalf("指纹库加载失败，请检查【fingerprint.txt】文件", err)
	} else {
		log.Log.Printf("成功加载HTTP指纹:[%d]条", n)
	}
	// First we parse the args
	flag.Parse()

	// Then we start the main process
	run.Start(&flag.Targets)

	// Finally we print the elapsed time
	elapsed := time.Since(startTime)
	log.Log.Printf("It costs %s", elapsed)
}
