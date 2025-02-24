package logger

import (
	"log"
	"runtime/debug"
)

func FatalIfError(err error) {
	fatalIf(err != nil, "fatal error: %v", err)
}

func fatalIf(cond bool, fmt string, args ...interface{}) {
	if !cond {
		return
	}
	debug.PrintStack()
	log.Fatalf(fmt, args...)
}

func Debugf(format string, args ...interface{}) {
	debug.PrintStack()
	log.Printf(format, args...)
}

func Infof(format string, args ...interface{}) {
	log.Printf(format, args...)
}

func Errorf(format string, args ...interface{}) {
	log.Println("--- [ERROR] ---")
	debug.PrintStack()
	log.Printf(format, args...)
	log.Println("----------------")
}
