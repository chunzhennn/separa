package main

import (
	"separa/common/flag"
)

func main() {
	Parse()

	if flag.Command.Cloud.Server {
		Server(flag.Command.Cloud.Port)
	} else {
		Node(flag.Command.Cloud.Port)
	}
}
