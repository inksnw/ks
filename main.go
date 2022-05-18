package main

import (
	"github.com/ks/rest"
	"github.com/phuslu/log"
	"os"
)

func initLog() {
	if !log.IsTerminal(os.Stderr.Fd()) {
		return
	}
	log.DefaultLogger = log.Logger{
		TimeFormat: "15:04:05",
		Caller:     1,
		Writer: &log.ConsoleWriter{
			ColorOutput:    true,
			QuoteString:    true,
			EndWithMessage: true,
		},
	}
}

func main() {
	initLog()
	rest.InitRestApi()
}
