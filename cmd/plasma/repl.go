package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/shoriwe/gplasma/pkg/compiler"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"os/signal"
	"strings"
)

var (
	input            = color.GreenString(">>>")
	waitingMoreInput = color.YellowString("...")
)

func splitFunc() bufio.SplitFunc {
	return func(data []byte, atEOF bool) (advance int, token []byte, err error) {
		_, compileError := compiler.Compile(string(data))
		if compileError == nil {
			return len(data), data, nil
		}
		if strings.Contains(compileError.Error(), "never ended") {
			_, _ = fmt.Fprint(os.Stdout, waitingMoreInput)
			return 0, nil, nil
		}
		return 0, nil, compileError
	}
}

func repl() {
	defer func() {
		err := recover()
		if err != nil {
			os.Exit(1)
		}
	}()
	plasma := vm.NewVM(os.Stdin, os.Stdout, os.Stderr)
	plasma.Load("exit", func(plasma *vm.Plasma) *vm.Value {
		return plasma.NewBuiltInFunction(plasma.Symbols(),
			func(argument ...*vm.Value) (*vm.Value, error) {
				if len(argument) == 0 {
					os.Exit(0)
				} else {
					os.Exit(int(argument[0].Int()))
				}
				return plasma.None(), nil
			})
	})
	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(splitFunc())
	onControlC := make(chan os.Signal, 1)
	signal.Notify(onControlC, os.Interrupt)
	for {
		_, _ = fmt.Fprint(os.Stdout, input)
		if !scanner.Scan() {
			onError("REPL", scanner.Err())
			continue
		}
		resultChannel, errorChannel, stopChannel := plasma.ExecuteString(scanner.Text())
		select {
		case <-onControlC:
			stopChannel <- struct{}{}
			onError("REPL", "Keyboard interruption")
		case executeError := <-errorChannel:
			result := <-resultChannel
			if executeError != nil {
				onError("REPL", executeError)
			} else if result != nil {
				if result != plasma.None() {
					_, _ = fmt.Fprintf(os.Stdout, "%s\n", result.String())
				}
			}
		}
	}
}
