package main

import (
	"bufio"
	"encoding/json"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"github.com/otiai10/copy"
	"github.com/shoriwe/gplasma"
	"github.com/shoriwe/gplasma/pkg/compiler"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/std/features/importlib"
	"github.com/shoriwe/gplasma/pkg/std/modules/base64"
	json2 "github.com/shoriwe/gplasma/pkg/std/modules/json"
	"github.com/shoriwe/gplasma/pkg/std/modules/regex"
	"github.com/shoriwe/gplasma/pkg/vm"
	"io"
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
	virtualMachine          *gplasma.VirtualMachine
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
	fmt.Printf("%s [MODE] [FLAG [FLAG [FLAG]]] [PROGRAM [PROGRAM [PROGRAM]]]\n", color.BlueString("%s", os.Args[0]))
	fmt.Printf("\n[%s] Notes\n", color.BlueString("+"))
	fmt.Printf("\t%s No %s arguments will spawn a %s\n", color.BlueString("-"), color.YellowString("PROGRAM"), color.YellowString("REPL"))
	fmt.Printf("\n[%s] Flags\n", color.BlueString("+"))
	for option, information := range flagOptions {
		fmt.Printf("\t%s, %s\t\t%s\n", color.RedString("%s", information.extra), color.RedString("%s", option), information.description)
	}
	fmt.Printf("\n[%s] Modes\n", color.BlueString("+"))
	fmt.Printf("\t%s\t\t%s\n", color.RedString("module"), "tool to install, uninstall and initialize modules")
	fmt.Printf("\n[%s] Environment Variables\n", color.BlueString("+"))
	for option, information := range envOptions {
		fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("%s", option), color.YellowString("%s", information.extra), information.description)
	}
	fmt.Println()
	os.Exit(0)
}

// Setup the vm based on the options
func setupVm() {
	virtualMachine = gplasma.NewVirtualMachine()
	currentDir, err := os.Getwd()
	if err != nil {
		currentDir = "."
	}
	importSystem := importlib.NewImporter()
	// Load Default modules to use with the VM
	importSystem.LoadModule(regex.Regex)
	importSystem.LoadModule(json2.JSON)
	importSystem.LoadModule(base64.Base64)
	//
	virtualMachine.LoadFeature(
		importSystem.Result(
			importlib.NewRealFileSystem(sitePackagesPath),
			importlib.NewRealFileSystem(currentDir),
		),
	)
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
	flagOptions = map[string]configOption{
		"--help": {
			extra:       "-h",
			description: "Show this help message",
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

func compileCode(scanner *bufio.Scanner) (*vm.Bytecode, *errors.Error) {
	scanner.Scan()
	sourceCode := scanner.Text()
	bytecode, compilationError := compiler.Compile(reader.NewStringReader(sourceCode))
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
									return compiler.Compile(reader.NewStringReader(sourceCode))
								}
							} else {
								return compiler.Compile(reader.NewStringReader(sourceCode))
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
						return compiler.Compile(reader.NewStringReader(sourceCode))
					}
				}
			}
			return nil, compilationError
		}
	}
	return bytecode, compilationError
}

func repl() {
	stdinScanner := bufio.NewScanner(os.Stdin)
	context := virtualMachine.NewContext()
	for {
		fmt.Print(color.GreenString(">>>"))
		bytecode, compilationError := compileCode(stdinScanner)
		if compilationError != nil {
			fmt.Printf("[%s] %s: %s\n", color.RedString("-"), compilationError.Type(), compilationError.Message())
			continue
		}

		result, success := virtualMachine.Execute(context, bytecode)
		if !success {
			fmt.Printf("[%s] %s: %s\n", color.RedString("-"), result.TypeName(), result.String)
			continue
		}
		if result.TypeName() == vm.NoneName {
			continue
		}
		resultToString, getError := result.Get(virtualMachine.Plasma, context, vm.ToString)
		if getError != nil {
			fmt.Println(result)
			continue
		}
		if !resultToString.IsTypeById(vm.FunctionId) {
			fmt.Println(result)
			continue
		}
		resultString, successToString := virtualMachine.CallFunction(
			context, resultToString,
		)
		if !successToString {
			fmt.Println(result)
			continue
		}
		fmt.Println(resultString.String)
	}
}

func program() {
	for _, filePath := range files {
		fileHandler, openError := os.Open(filePath)
		if openError != nil {
			_, _ = fmt.Fprintf(os.Stderr, openError.Error())
			os.Exit(1)
		}
		content, readingError := io.ReadAll(fileHandler)
		if readingError != nil {
			_, _ = fmt.Fprintf(os.Stderr, readingError.Error())
			os.Exit(1)
		}
		result, success := virtualMachine.ExecuteMain(string(content))
		if !success {
			_, _ = fmt.Fprintf(os.Stderr, "[%s] %s: %s\n", color.RedString("-"), result.TypeName(), result.String)
			os.Exit(1)
		}
	}
}

func helpModules() {
	fmt.Printf("%s module OPTION ARGUMENT\n", color.BlueString("%s", os.Args[0]))
	fmt.Printf("\n[%s] Options\n", color.BlueString("+"))
	fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("install"), color.YellowString("MODULE_PATH"), "Install a module in path")
	fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("uninstall"), color.YellowString("MODULE_PATH"), "Uninstall a module in path")
	fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("init"), color.YellowString("MODULE_PATH"), "initialize a new module")
	fmt.Printf("\n[%s] Environment Variables\n", color.BlueString("+"))
	for option, information := range envOptions {
		fmt.Printf("\t%s -> %s\t\t%s\n", color.RedString("%s", option), color.YellowString("%s", information.extra), information.description)
	}
	fmt.Println()
	os.Exit(0)
}

func installModuleFromDisk(source string) {
	_, err := os.Stat(source)
	if err != nil {
		if os.IsExist(err) {
			_, _ = fmt.Printf("[%s] No sufficient permissions to open the path %s\n", color.RedString("-"), color.RedString(source))
			os.Exit(1)
		} else if os.IsNotExist(err) {
			_, _ = fmt.Printf("[%s] Path %s doesn't exists\n", color.RedString("-"), color.RedString(source))
			os.Exit(1)
		} else {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		}
	}
	settingsFile, openError := os.Open(filepath.Join(source, "settings.json"))
	if openError != nil {
		_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), openError.Error())
		os.Exit(1)
	}
	content, readingError := io.ReadAll(settingsFile)
	if readingError != nil {
		_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), readingError.Error())
		os.Exit(1)
	}
	var settings importlib.Settings
	unmarshalError := json.Unmarshal(content, &settings)
	if unmarshalError != nil {
		_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), unmarshalError.Error())
		os.Exit(1)
	}
	modulePath := filepath.Join(sitePackagesPath, settings.Name)
	_, err = os.Stat(modulePath)
	if err != nil {
		if os.IsExist(err) {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		} else if os.IsNotExist(err) {
			creationError := os.Mkdir(modulePath, 0755)
			if creationError != nil {
				_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), creationError.Error())
				os.Exit(1)
			}
		} else {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		}
	}
	moduleVersionPath := filepath.Join(modulePath, settings.Version)
	_, err = os.Stat(moduleVersionPath)
	if err != nil {
		if os.IsExist(err) {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		} else if os.IsNotExist(err) {
			creationError := os.Mkdir(moduleVersionPath, 0755)
			if creationError != nil {
				_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), creationError.Error())
				os.Exit(1)
			}
		} else {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), err.Error())
			os.Exit(1)
		}
	} else {
		deletionError := os.RemoveAll(moduleVersionPath)
		if deletionError != nil {
			_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), deletionError.Error())
			os.Exit(1)
		}
	}
	for _, dependency := range settings.Dependencies {
		installModule(dependency)
	}
	copyError := copy.Copy(source, moduleVersionPath)
	if copyError != nil {
		_, _ = fmt.Printf("[%s] %s\n", color.RedString("-"), copyError.Error())
		os.Exit(1)
	}
}

func installModule(module string) {
	if strings.Index(module, "github.com/") == 0 {
		// ToDo: create the part of the tool that handle github repositories
	} else {
		installModuleFromDisk(module)
	}
}

func moduleInstallation() {
	if len(os.Args) != 4 {
		fmt.Printf("%s module uninstall MODULE_NAME[@MODULE_VERSION]", os.Args[0])
		os.Exit(0)
	}
	module := os.Args[3]
	installModule(module)
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
		_, _ = fmt.Printf(
			"[%s] Invalid nomenclature of %s",
			color.RedString("-"),
			color.RedString("MODULE@VERSION"),
		)
		os.Exit(0)
	}
	moduleName := splitModule[0]
	if len(splitModule) == 2 {
		version = splitModule[1]
	}
	modulePath := filepath.Join(sitePackagesPath, moduleName)
	_, err := os.Stat(modulePath)
	if err != nil {
		if os.IsNotExist(err) {
			_, _ = fmt.Printf("[%s] No module with name %s", color.RedString("-"), color.RedString(moduleName))
			os.Exit(0)
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
			_, _ = fmt.Printf("[%s] No module with name %s and version %s", color.RedString("-"), color.RedString(moduleName), color.YellowString(version))
			os.Exit(0)
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

func initializeModule() {
	if len(os.Args) != 4 {
		fmt.Printf("%s module uninstall MODULE_NAME[@MODULE_VERSION]", os.Args[0])
		os.Exit(0)
	}
	location, moduleName := filepath.Split(os.Args[3])
	location = filepath.Clean(location)
	_, err := os.Stat(os.Args[3])
	if err != nil {
		if os.IsNotExist(err) {
			folderCreationError := os.Mkdir(os.Args[3], 0755)
			if folderCreationError != nil {
				_, _ = fmt.Printf("[%s] %s", color.RedString("-"), folderCreationError.Error())
				os.Exit(0)
			}
		} else {
			_, _ = fmt.Printf("[%s] %s", color.RedString("-"), err.Error())
			os.Exit(0)
		}
	}
	settingsFile, openError := os.OpenFile(filepath.Join(os.Args[3], "settings.json"), os.O_CREATE|os.O_RDWR, 0755)
	if openError != nil {
		_, _ = fmt.Printf("[%s] %s", color.RedString("-"), openError.Error())
		os.Exit(0)
	}
	settings := importlib.Settings{
		Name:         moduleName,
		Version:      "0.0.0",
		Resources:    "",
		EntryScript:  "",
		Dependencies: []string{},
	}
	jsonSettings, marshalError := json.Marshal(settings)
	if marshalError != nil {
		_, _ = fmt.Printf("[%s] %s", color.RedString("-"), marshalError.Error())
		os.Exit(0)
	}
	_, writeError := settingsFile.Write(jsonSettings)
	if writeError != nil {
		_, _ = fmt.Printf("[%s] %s", color.RedString("-"), writeError.Error())
		os.Exit(0)
	}
	closeError := settingsFile.Close()
	if closeError != nil {
		_, _ = fmt.Printf("[%s] %s", color.RedString("-"), closeError.Error())
		os.Exit(0)
	}
}

func modules() {
	if len(os.Args) == 2 {
		helpModules()
		os.Exit(0)
	}
	switch os.Args[2] {
	case "install":
		moduleInstallation()
	case "uninstall":
		uninstallModule()
	case "init":
		initializeModule()
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
