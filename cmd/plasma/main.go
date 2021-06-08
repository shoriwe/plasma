package main

import (
	"fmt"
	"github.com/fatih/color"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
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
	// fmt.Printf("\t%s No %s arguments will spawn a %s\n", color.BlueString("-"), color.YellowString("PROGRAM"), color.YellowString("REPL"))
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
	help()
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
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s", color.RedString("-"), executionError.GetString())
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
