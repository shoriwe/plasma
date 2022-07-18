package main

import (
	"fmt"
	"github.com/fatih/color"
	"os"
)

func onError(file string, a ...any) {
	_, _ = fmt.Fprint(os.Stderr, color.RedString("%s: %s\n", file, fmt.Sprint(a)))
}
