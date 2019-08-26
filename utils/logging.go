package utils

import (
	"flag"
	"io/ioutil"
	"log"
)

// Enable logging if -l (-log) flag is set.
func InitLogging() {
	logEnabled := flag.Bool("l", false, "enables logging")
	logEnabledLong := flag.Bool("log", false, "enables logging")
	flag.Parse()
	if !*logEnabled && !*logEnabledLong {
		log.SetFlags(0)
		log.SetOutput(ioutil.Discard)
	}
}
