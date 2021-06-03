package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/shoriwe/gruby/pkg/compiler/plasma"
	"github.com/shoriwe/gruby/pkg/reader"
	"github.com/shoriwe/gruby/pkg/vm"
	"os"
)

const (
	REPL = iota
	Program

	QuickHelp           = "-h"
	Help                = "--help"
	QuickDisableImports = "-I"
	DisableImports      = "--no-imports"

	NoColor = "NoColor"
)

var (
	files          []string
	virtualMachine *vm.Plasma
	flagOptions    map[string]configOption
	envOptions     map[string]configOption
	mode           = REPL
)

type configOption struct {
	extra       string
	description string
	onFound     func()
}

// Environment Variables functions
func noColor() {
	value := os.Getenv(NoColor)
	if value == "TRUE" {
		color.NoColor = true
	} else {
		color.NoColor = false
	}
}

// Flags functions
func help() {
	fmt.Printf("%s [FLAG [FLAG [FLAG]]] [PROGRAM [PROGRAM [PROGRAM]]]\n", color.BlueString("%s", os.Args[0]))
	fmt.Printf("\n[%s] Notes\n", color.BlueString("+"))
	fmt.Printf("\t%s No %s arguments will spawn a %s\n", color.BlueString("-"), color.YellowString("PROGRAM"), color.YellowString("REPL"))
	fmt.Printf("\n[%s] Flags\n", color.BlueString("+"))
	for option, information := range flagOptions {
		fmt.Printf("\t%s, %s\t\t%s\n", color.RedString("%s", information.extra), color.RedString("%s", option), information.description)
	}
	fmt.Printf("\n[%s] Environment Variables\n", color.BlueString("+"))
	for option, information := range envOptions {
		fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("%s", option), color.YellowString("%s", information.extra), information.description)
	}
	fmt.Println()
	os.Exit(0)
}

// Setup the vm based on the options
func setupVm() {
	if len(files) > 0 {
		mode = Program
	}
	virtualMachine = vm.NewPlasmaVM(os.Stdin, os.Stdout, os.Stderr)
	// Setup from here the other flags
}

func init() {
	flagOptions = map[string]configOption{
		Help: {
			extra:       QuickHelp,
			description: "Show this help message",
			onFound:     help,
		},
	}
	envOptions = map[string]configOption{
		NoColor: {
			extra:       fmt.Sprintf("%s or %s", color.YellowString("TRUE"), color.YellowString("FALSE")),
			description: "Disable color printing for this CLI",
			onFound:     noColor,
		},
	}
	for _, information := range envOptions {
		information.onFound()
	}
	for index, argument := range os.Args[1:] {
		if argument[0] == '-' {
			switch argument {
			case QuickHelp, Help:
				help()
			case QuickDisableImports, DisableImports:
				break // ToDo: Implement me
			default:
				_, _ = fmt.Fprintf(os.Stderr, "Option %s -> %s Unknown", color.BlueString("%d", index+1), color.RedString("%s", argument))
				os.Exit(1)
			}
		} else {
			files = append(files, argument)
		}
	}
	setupVm()
}

func repl() {
	stdin := bufio.NewReader(os.Stdin)
	prefix := ">>>"
	virtualMachine.Initialize(make([]vm.Code, 0))
	var code string
	for ; ; {
		fmt.Print(color.BlueString(prefix))
		inputLine, readingError := stdin.ReadString('\n')
		if readingError != nil {
			panic(readingError)
		}
		if inputLine == "\n" {
			continue
		}
		code += inputLine
		compiler := plasma.NewCompiler(reader.NewStringReader(code), nil)
		compiledCode, compilationError := compiler.Compile()
		if compilationError != nil {
			// ToDo: Do something with the signal received
		}
		virtualMachine.PushCode(compiledCode)
		result, executionError := virtualMachine.Execute()
		if virtualMachine.MemoryStack.HasNext() {
			virtualMachine.PopObject()
		}
		// virtualMachine.PopCode()
		if executionError != nil {
			fmt.Printf("[%s] %s\n", color.RedString("-"), color.RedString(executionError.String()))
			continue
		} else if result.TypeName() == vm.NoneName {
			continue
		}
		toString, getError := result.Get(vm.ToString)
		if getError != nil {
			fmt.Printf("\n[%s] Object does not implement method %s\n", color.RedString("-"), color.YellowString(vm.ToString))
			continue
		}
		if _, ok := toString.(*vm.Function); !ok {
			fmt.Printf("\n[%s] Object %s is not a function with 0 arguments\n", color.RedString("-"), color.YellowString(vm.ToString))
			continue
		}
		stringResult, callError := vm.CallFunction(toString.(*vm.Function), virtualMachine, result.SymbolTable())
		if callError != nil {
			fmt.Printf("\n[%s] \"%s\"\n", color.RedString("-"), color.RedString(callError.String()))
			continue
		}
		fmt.Println(stringResult.GetString())
	}
}

func program() {
	for _, file := range files {
		fileHandler, readingError := os.Open(file)
		if readingError != nil {
			panic(readingError)
		}
		compiler := plasma.NewCompiler(reader.NewStringReaderFromFile(fileHandler),
			map[uint8]uint8{
				plasma.PopRawExpressions: plasma.PopRawExpressions,
			},
		)
		code, compilationError := compiler.Compile()
		if compilationError != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s", color.RedString("-"), compilationError.String())
			os.Exit(1)
		}
		/*
			ToDo: Do intermediate stuff with other flags
		*/
		virtualMachine.InitializeByteCode(code)
		_, executionError := virtualMachine.Execute()
		if executionError != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s", color.RedString("-"), executionError.String())
			os.Exit(1)
		}
	}
}

func main() {
	switch mode {
	case REPL:
		repl()
	case Program:
		program()
	}
}
