package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/std/importlib"
	"github.com/shoriwe/gplasma/pkg/vm"
	"os"
	"path/filepath"
	"strings"
)

const (
	NoColor      = "NoColor"
	SitePackages = "SitePackages"
)

var (
	files                   []string
	virtualMachine          *vm.Plasma
	flagOptions             map[string]configOption
	envOptions              map[string]configOption
	defaultSitePackagesPath = true
	sitePackagesPath        = "site-packages"
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

func sitePackages() {
	value := os.Getenv(SitePackages)
	if value != "" {
		sitePackagesPath = os.Getenv(SitePackages)
		defaultSitePackagesPath = false
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
	virtualMachine = vm.NewPlasmaVM(os.Stdin, os.Stdout, os.Stderr)
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}
	virtualMachine.LoadBuiltInSymbols(
		importlib.NewImporter(
			importlib.NewRealFileSystem(sitePackagesPath),
			importlib.NewRealFileSystem(currentDir),
		),
	)
	// Setup from here the other flags
}

func init() {
	envOptions = map[string]configOption{
		NoColor: {
			extra:       fmt.Sprintf("%s or %s", color.YellowString("TRUE"), color.YellowString("FALSE")),
			description: "Disable color printing for this CLI",
			onFound:     noColor,
		},
		SitePackages: {
			extra:       fmt.Sprintf("%s", color.YellowString("PATH")),
			description: fmt.Sprintf("This is the path to the Site-Packages of the running VM; Default is %s", color.BlueString("PATH/TO/PLASMA/EXECUTABLE/%s", sitePackagesPath)),
			onFound:     sitePackages,
		},
	}
	for _, information := range envOptions {
		information.onFound()
	}
	if defaultSitePackagesPath {
		exePath, err := os.Executable()
		if err != nil {
			panic(err)
		}
		directory, _ := filepath.Split(exePath)
		sitePackagesPath = filepath.Join(directory, sitePackagesPath)
	}
	if _, err := os.Stat(sitePackagesPath); err != nil {
		if os.IsNotExist(err) {
			dirCreationError := os.Mkdir(sitePackagesPath, 0755)
			if dirCreationError != nil {
				_, _ = fmt.Fprintf(os.Stderr, "[%s] %s\n", color.RedString("-"), dirCreationError.Error())
				os.Exit(1)
			}
		} else {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		}
	}
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
			fmt.Printf("[%s] %s: %s\n", color.RedString("-"), compilationError.Type(), compilationError.Message())
			continue
		}
		virtualMachine.PushBytecode(vm.NewBytecodeFromArray(bytecode))
		result, executionError := virtualMachine.Execute()
		virtualMachine.PopBytecode()
		if executionError != nil {
			fmt.Printf("[%s] %s: %s\n", color.RedString("-"), executionError.TypeName(), executionError.GetString())
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
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s\n", color.RedString("-"), compilationError.String())
			os.Exit(1)
		}
		/*
			ToDo: Do intermediate stuff with other flags
		*/
		virtualMachine.InitializeBytecode(code)
		_, executionError := virtualMachine.Execute()
		if executionError != nil {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", color.RedString("-"), executionError.TypeName(), executionError.GetString())
			os.Exit(1)
		}
	}
}

func helpModules() {
	fmt.Printf("%s module OPTION ARGUMENT\n", color.BlueString("%s", os.Args[0]))
	fmt.Printf("\n[%s] Options\n", color.BlueString("+"))
	fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("install"), color.YellowString("MODULE_PATH"), "Install a module in path")
	fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("uninstall"), color.YellowString("MODULE_PATH"), "Uninstall a module in path")
	fmt.Printf("\n[%s] Environment Variables\n", color.BlueString("+"))
	for option, information := range envOptions {
		fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("%s", option), color.YellowString("%s", information.extra), information.description)
	}
	fmt.Println()
	os.Exit(0)
}

func installModule() {

}

func uninstallModule() {
	if len(os.Args) != 4 {
		fmt.Printf("%s module uninstall MODULE_NAME[@MODULE_VERSION]", os.Args[0])
		os.Exit(0)
	}
	module := os.Args[3]
	splitModule := strings.Split(module, "@")
	version := "all"
	if len(splitModule) > 2 {
		panic("Invalid nomenclature of MODULE@VERSION")
	}
	moduleName := splitModule[0]
	if len(splitModule) == 2 {
		version = splitModule[1]
	}
	modulePath := filepath.Join(sitePackagesPath, moduleName)
	_, err := os.Stat(modulePath)
	if err != nil {
		if os.IsNotExist(err) {
			panic("No module with name " + moduleName + " installer")
		} else {
			panic(err)
		}
	}
	if version == "all" {
		removeError := os.RemoveAll(modulePath)
		if removeError != nil {
			panic(removeError)
		}
		os.Exit(0)
	}
	modulePath = filepath.Join(modulePath, version)
	_, err = os.Stat(modulePath)
	if err != nil {
		if os.IsNotExist(err) {
			panic("No module with name " + moduleName + " installer")
		} else {
			panic(err)
		}
	}
	removeError := os.RemoveAll(modulePath)
	if removeError != nil {
		panic(removeError)
	}
	os.Exit(0)
}

func modules() {
	if len(os.Args) == 2 {
		helpModules()
		os.Exit(0)
	}
	switch os.Args[2] {
	case "install":
		installModule()
	case "uninstall":
		uninstallModule()
	default:
		helpModules()
		os.Exit(1)
	}
}

func execution() {
	if len(os.Args) == 1 {
		setupVm()
		repl()
		os.Exit(0)
	}
	for _, arg := range os.Args {
		if arg == "-h" || arg == "--help" {
			help()
			os.Exit(0)
		}
	}
	plasmaFlagSet := flag.NewFlagSet("plasma", flag.ExitOnError)
	parsingError := plasmaFlagSet.Parse(os.Args[1:])
	if parsingError != nil {
		panic(parsingError)
	}

	if len(plasmaFlagSet.Args()) == 0 {
		setupVm()
		repl()
	} else {
		files = plasmaFlagSet.Args()
		setupVm()
		program()
	}
}

func main() {
	if len(os.Args) < 2 {
		execution()
		os.Exit(0)
	}
	switch os.Args[1] {
	case "module":
		modules()
	default:
		execution()
	}
}
