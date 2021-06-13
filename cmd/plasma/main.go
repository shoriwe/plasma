package main

import (
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"strings"
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

func compileCode(scanner *bufio.Scanner) ([]vm.Code, *errors.Error) {
	scanner.Scan()
	sourceCode := scanner.Text()
	compiler := plasma.NewCompiler(reader.NewStringReader(sourceCode), plasma.Options{})
	bytecode, compilationError := compiler.CompileToArray()
	if compilationError != nil {
		switch compilationError.Type() {
		case errors.LexingError:
			if strings.Contains(compilationError.Message(), "string never closed") {
				for {
					fmt.Print(color.BlueString("..."))
					scanner.Scan()
					rest := scanner.Text()
					restLength := len(rest)
					sourceCode += "\n" + rest
					if restLength >= 1 {
						if sourceCode[0] == rest[restLength-1] {
							if restLength-2 >= 0 {
								if rest[restLength-2] != '\\' {
									return plasma.NewCompiler(reader.NewStringReader(sourceCode), plasma.Options{}).CompileToArray()
								}
							} else {
								return plasma.NewCompiler(reader.NewStringReader(sourceCode), plasma.Options{}).CompileToArray()
							}
						}
					}
				}
			}
			return nil, compilationError
		case errors.SyntaxError:
			if strings.Contains(compilationError.Message(), "invalid definition of") {
				for {
					fmt.Print(color.BlueString("..."))
					scanner.Scan()
					rest := scanner.Text()
					sourceCode += "\n" + rest
					if rest == "end" {
						return plasma.NewCompiler(reader.NewStringReader(sourceCode), plasma.Options{}).CompileToArray()
					}
				}
			}
			return nil, compilationError
		}
	}
	return bytecode, compilationError
}

func repl() {
	virtualMachine.InitializeBytecode(nil)
	stdinScanner := bufio.NewScanner(os.Stdin)
	for {
		fmt.Print(color.GreenString(">>>"))
		bytecode, compilationError := compileCode(stdinScanner)
		if compilationError != nil {
			fmt.Printf("[%s] %s: %s", color.RedString("-"), compilationError.Type(), compilationError.Message())
			continue
		}
		virtualMachine.PushBytecode(vm.NewBytecodeFromArray(bytecode))
		result, executionError := virtualMachine.Execute()
		virtualMachine.PopBytecode()
		if executionError != nil {
			fmt.Printf("[%s] %s: %s", color.RedString("-"), executionError.TypeName(), executionError.GetString())
			continue
		}
		if result.TypeName() == vm.NoneName {
			continue
		}
		resultToString, getError := result.Get(vm.ToString)
		if getError != nil {
			fmt.Println(result)
			continue
		}
		if _, ok := resultToString.(*vm.Function); !ok {
			fmt.Println(result)
			continue
		}
		resultString, callError := virtualMachine.CallFunction(resultToString.(*vm.Function), virtualMachine.PeekSymbolTable())
		if callError != nil {
			fmt.Println(result)
			continue
		}
		fmt.Println(resultString.GetString())
	}
}

func program() {
	virtualMachine.InitializeBytecode(nil)
	for _, file := range files {
		fileHandler, readingError := os.Open(file)
		if readingError != nil {
			panic(readingError)
		}
		compiler := plasma.NewCompiler(reader.NewStringReaderFromFile(fileHandler),
			plasma.Options{
				Debug:             false,
				PopRawExpressions: true,
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
		virtualMachine.InitializeBytecode(code)
		_, executionError := virtualMachine.Execute()
		if executionError != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s: %s", color.RedString("-"), executionError.TypeName(), executionError.GetString())
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
