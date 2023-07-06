package flag

import (
	"fmt"
	"io/ioutil"
	"os"
	"separa/common/log"
	"strings"

	"github.com/alecthomas/kong"
)

var Command struct {
	// Debug Mode
	Debug bool `help:"Debug mode to get more detail output."`

	// Target CIDR
	Target string `help:"Target to scan, supports CIDR." name:"target" short:"t"`

	// CIDRs in a file split by '\n'
	TargetFile string `help:"Target file to scan, split each target line by line with '\\n'." name:"target-file" short:"f"`

	// Config file
	ConfigFile string `help:"Config file to load, default is config.yaml in current dir." name:"config-file" type:"path" short:"c" default:"config.yaml"`

	// Output file
	OutputFile string `help:"Output file to save, default is output.json in current dir." name:"output-file" type:"path" short:"o" default:"output.json"`

	// Port to scan
	Port string `help:"Port to scan, default is TOP 1000. you can use ',' to split or '-' to range, like '80,443,22' or '1-65535'" name:"port" short:"p"`
}

func CheckTarget() {
	if Command.Target == "" && Command.TargetFile == "" {
		fmt.Println("Missing target! use '--help/-h' for help.")
		os.Exit(0)
	}
}

var Targets []string

// Parse Args to scan
func Parse() (ctx *kong.Context) {
	ctx = kong.Parse(
		&Command,
		kong.Name("separa"),
		kong.Description("A simple scanner for Web Security"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
	)

	// Check if target is empty
	CheckTarget()

	log.Log.Printf("Args Load: %+v", Command)

	if Command.Target != "" {
		Targets = append(Targets, Command.Target)

	}
	if Command.TargetFile != "" {
		content, err := ioutil.ReadFile(Command.TargetFile)
		if err != nil {
			log.Log.Fatalf("Read target file error: %s", err)
			return ctx
		}
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			Targets = append(Targets, strings.TrimSpace(line))
		}
	}

	defer log.Log.Printf("%d CIDR Load", len(Targets))

	return ctx
}
