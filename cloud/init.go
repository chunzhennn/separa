package main

import (
	"fmt"
	"os"
	"separa/common/flag"

	"github.com/alecthomas/kong"
)

// Parse Args to init
func Parse() {
	ctx := kong.Parse(
		&flag.Command,
		kong.Name("separa"),
		kong.Description("A simple scanner for Web Security"),
		kong.UsageOnError(),
		kong.ConfigureHelp(kong.HelpOptions{
			Compact: true,
			Summary: false,
		}),
	)

	if ctx.Command() != "cloud" {
		fmt.Println("mutipara cann't use 'cloud' command!")
		os.Exit(0)
	}
}
