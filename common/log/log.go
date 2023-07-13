package log

import (
	"log"
	"separa/common/flag"
	// "go.uber.org/zap"
)

var Log = log.Default()

func Out(format string, v ...any) {
	Log.Printf("\x1B[32m[+]\x1B[0m "+format+"\n", v...)
}

func Dbg(format string, v ...any) {
	if flag.Command.Debug {
		Log.Printf("\x1B[33m[*]\x1B[0m "+format+"\n", v...)
	}
}

func Err(format string, v ...any) {
	Log.Printf("\x1B[31m[x]\x1B[0m "+format+"\n", v...)
}
