package flag

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/alecthomas/kong"
)

var Command struct {
	Scan struct {
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

		// Delay
		Delay int `help:"Delay between each request " name:"delay" short:"d" default:"5"`

		// Top N ports to scan
		Top int `help:"Top N ports to scan, default is 1000" name:"top" short:"n" default:"1000"`
	} `cmd:""`

	Cloud struct {
		// Server Mode
		Server bool `help:"Server mode to add node." default:"false"`
		// Port
		Port string `help:"Port to listen, default is 8080" default:"8080"`
	} `cmd:""`
}

func CheckTarget() {
	if Command.Scan.Target == "" && Command.Scan.TargetFile == "" {
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

	if ctx.Command() != "scan" {
		fmt.Println("separa cann't use 'scan' command!")
		os.Exit(0)
	}

	// Check if target is empty
	CheckTarget()

	fmt.Printf("Args Load: %+v\n", Command)

	if Command.Scan.Target != "" {
		Targets = append(Targets, Command.Scan.Target)

	}
	if Command.Scan.TargetFile != "" {
		content, err := ioutil.ReadFile(Command.Scan.TargetFile)
		if err != nil {
			fmt.Printf("Read target file error: %s\n", err)
			return ctx
		}
		lines := strings.Split(string(content), "\n")
		for _, line := range lines {
			Targets = append(Targets, strings.TrimSpace(line))
		}
	}

	defer fmt.Printf("%d CIDR Load\n", len(Targets))

	return ctx
}
