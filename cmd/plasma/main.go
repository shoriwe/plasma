package main

import (
	"os"
)

func main() {
	if len(os.Args) == 1 {
		repl()
	} else if len(os.Args) > 1 {
		executeFiles()
	}
}
