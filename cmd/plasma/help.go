package main

import (
	"fmt"
	"os"
)

const helpMessage = `Usage: %s [FILE [FILE [FILE [...]]]]

Zero arguments will start the REPL'`

func help() {
	fmt.Printf(helpMessage, os.Args[0])
	os.Exit(0)
}
