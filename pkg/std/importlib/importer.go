package importlib

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"github.com/shoriwe/gplasma/pkg/compiler/plasma"
	"github.com/shoriwe/gplasma/pkg/errors"
	"github.com/shoriwe/gplasma/pkg/reader"
	"github.com/shoriwe/gplasma/pkg/vm"
	"io"
	"path/filepath"
	"strings"
)

type context struct {
	moduleName  string
	version     string
	resources   string
	entryScript string
	root        string
}

type Settings struct {
	Name         string
	Version      string
	Resources    string
	EntryScript  string
	Dependencies []string
}

const (
	ResourceReader          = "ResourceReader"
	ResourceNotFoundError   = "ResourceNotFoundError"
	NotInsideModuleError    = "NotInsideModuleError"
	ScriptNotFoundError     = "ScriptNotFoundError"
	CompilationError        = "CompilationError"
	ModuleNotFoundError     = "ModuleNotFoundError"
	ChangeDirectoryError    = "ChangeDirectoryError"
	ModuleNomenclatureError = "ModuleNomenclatureError"
	NoVersionFoundError     = "NoVersionFoundError"
)

func getScriptHash(r io.ReadSeekCloser) (string, error) {
	table := md5.New()
	for {
		chunk := make([]byte, 1000000)
		length, readError := r.Read(chunk)
		if readError != nil {
			if readError == io.EOF {
				break
			}
			return "", readError
		}
		table.Write(chunk[:length])
	}
	_, seekError := r.Seek(0, io.SeekStart)
	if seekError != nil {
		return "", seekError
	}
	return string(table.Sum(nil)), nil
}

func resourceReaderInitialize(p *vm.Plasma, r io.ReadSeekCloser) vm.ConstructorCallBack {
	return func(object vm.Value) *vm.Object {
		object.Set("Read",
			p.NewFunction(false, object.SymbolTable(),
				vm.NewBuiltInClassFunction(object, 1,
					func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
						bytesToRead := arguments[0]
						if _, ok := bytesToRead.(*vm.Integer); !ok {
							return p.NewInvalidTypeError(bytesToRead.TypeName(), vm.IntegerName), nil
						}
						bytes := make([]byte, bytesToRead.GetInteger().Int64())
						numberOfBytes, readError := r.Read(bytes)
						if readError != nil {
							if readError == io.EOF {
								return p.NewNone(), nil
							} else {
								return nil, p.NewGoRuntimeError(readError)
							}
						}
						bytes = bytes[:numberOfBytes]
						return p.NewBytes(false, p.PeekSymbolTable(), bytes), nil
					},
				),
			),
		)
		object.Set("Seek",
			p.NewFunction(false, object.SymbolTable(),
				vm.NewBuiltInClassFunction(object, 1,
					func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
						seek := arguments[0]
						if _, ok := seek.(*vm.Integer); !ok {
							return p.NewInvalidTypeError(seek.TypeName(), vm.IntegerName), nil
						}
						_, seekError := r.Seek(seek.GetInteger().Int64(), io.SeekStart)
						if seekError != nil {
							return nil, p.NewGoRuntimeError(seekError)
						}
						return p.NewNone(), nil
					},
				),
			),
		)
		object.Set("Close",
			p.NewFunction(false, object.SymbolTable(),
				vm.NewBuiltInClassFunction(object, 0,
					func(self vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
						closeError := r.Close()
						if closeError != nil {
							return nil, p.NewGoRuntimeError(closeError)
						}
						return p.NewNone(), nil
					},
				),
			),
		)
		return nil
	}
}

func newResourceReader(p *vm.Plasma, r io.ReadSeekCloser) vm.Value {
	resourceReader := p.NewObject(false, ResourceReader, nil, p.PeekSymbolTable())
	resourceReaderInitialize(p, r)(resourceReader)
	return resourceReader
}

func newResourceNotFoundError(p *vm.Plasma, path string) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(ResourceNotFoundError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), path),
	)
	return result.(*vm.Object)
}

func newNotInsideModuleError(p *vm.Plasma) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(NotInsideModuleError))
	p.ForceInitialization(result)
	return result.(*vm.Object)
}

func newScriptNotFoundError(p *vm.Plasma, path string) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(ScriptNotFoundError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), path),
	)
	return result.(*vm.Object)
}

func newCompilationError(p *vm.Plasma, compilationError *errors.Error) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(CompilationError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), compilationError.Message()),
	)
	return result.(*vm.Object)
}

func newModuleNotFoundError(p *vm.Plasma, moduleName string) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(ModuleNotFoundError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), moduleName),
	)
	return result.(*vm.Object)
}

func newChangeDirectoryError(p *vm.Plasma, compilationError *errors.Error) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(ChangeDirectoryError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), compilationError.Message()),
	)
	return result.(*vm.Object)
}

func newModuleNomenclatureError(p *vm.Plasma, moduleName string) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(ModuleNomenclatureError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), moduleName),
	)
	return result.(*vm.Object)
}

func newNoVersionFoundError(p *vm.Plasma, moduleName string, version string) *vm.Object {
	result := p.ForceConstruction(p.ForceMasterGetAny(NoVersionFoundError))
	p.ForceInitialization(result,
		p.NewString(false, p.PeekSymbolTable(), moduleName),
		p.NewString(false, p.PeekSymbolTable(), version),
	)
	return result.(*vm.Object)
}

func (c *context) isSet() bool {
	return c.root != ""
}

func getResource(ctx *context, sitePackages FileSystem) vm.ObjectLoader {
	return func(p *vm.Plasma) vm.Value {
		return p.NewFunction(true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(_ vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					if !ctx.isSet() {
						return nil, newNotInsideModuleError(p)
					}

					resourcePathObject := arguments[0]
					if _, ok := resourcePathObject.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(resourcePathObject.TypeName(), vm.StringName)
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					resourcePath := resourcePathObject.GetString()
					if !sitePackages.ExistsRelative(resourcePath) {
						return nil, newResourceNotFoundError(p, resourcePath)
					}
					resourceHandler, openError := sitePackages.OpenRelative(resourcePath)
					if openError != nil {
						return nil, p.NewGoRuntimeError(openError)
					}
					return newResourceReader(p, resourceHandler), nil
				},
			),
		)
	}
}

func getResourcePath(ctx *context, sitePackages FileSystem) vm.ObjectLoader {
	return func(p *vm.Plasma) vm.Value {
		return p.NewFunction(true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					resource := arguments[0]
					if _, ok := resource.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(resource.TypeName(), vm.StringName)
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					if !sitePackages.ExistsRelative(resource.GetString()) {
						return nil, newResourceNotFoundError(p, resource.GetString())
					}
					resourcePath := filepath.Join(sitePackages.AbsolutePwd(), ctx.root, ctx.resources, resource.GetString())
					return p.NewString(false, p.PeekSymbolTable(), resourcePath), nil
				},
			),
		)
	}
}

func scriptImport(memory map[string]vm.Value, ctx *context, sitePackages FileSystem, pwd FileSystem) vm.ObjectLoader {
	return func(p *vm.Plasma) vm.Value {
		return p.NewFunction(true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					scriptPath := arguments[0]
					if _, ok := scriptPath.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(scriptPath.TypeName(), vm.StringName)
					}
					// First try to get the content of the file
					var (
						scriptFile io.ReadSeekCloser
						openError  error
					)
					// Check if the code is being ran in the context of a module
					if ctx.isSet() {
						// If it is, import the scriptFile relative the root of the module
						// Check if file exists
						if !sitePackages.ExistsRelative(scriptPath.GetString()) {
							return nil, newScriptNotFoundError(p, scriptPath.GetString())
						}
						scriptFile, openError = sitePackages.OpenRelative(scriptPath.GetString())

						oldLocation := sitePackages.RelativePwd()
						defer sitePackages.ChangeDirectoryFullPath(oldLocation)
						sitePackages.ChangeDirectoryToFileLocation(scriptPath.GetString())
					} else {
						// If not import the scriptFile from the immediate filesystem of the running scriptFile
						// Check if file exists
						if !pwd.ExistsRelative(scriptPath.GetString()) {
							return nil, newScriptNotFoundError(p, scriptPath.GetString())
						}
						scriptFile, openError = pwd.OpenRelative(scriptPath.GetString())
						oldLocation := pwd.RelativePwd()
						defer pwd.ChangeDirectoryFullPath(oldLocation)
						pwd.ChangeDirectoryToFileLocation(scriptPath.GetString())
					}
					if openError != nil {
						return nil, p.NewGoRuntimeError(openError)
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return nil, p.NewGoRuntimeError(hashingError)
					}
					if _, ok := memory[scriptHash]; ok {
						return memory[scriptHash], nil
					}
					// ToDo: Fix this, use a better file reader object
					scriptCode, compilationError := plasma.NewCompiler(reader.NewStringReaderFromFile(scriptFile),
						plasma.Options{
							Debug: false,
						},
					).Compile()
					if compilationError != nil {
						return nil, newCompilationError(p, compilationError)
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(false, p.PeekSymbolTable())
					memory[scriptHash] = script
					p.PushSymbolTable(script.SymbolTable())
					p.PushBytecode(scriptCode)
					_, executionError := p.Execute()
					if executionError != nil {
						return nil, executionError
					}
					p.PopBytecode()
					p.PopSymbolTable()
					// Return the initialized module object
					return script, nil
				},
			),
		)
	}
}

func moduleImport(memory map[string]vm.Value, ctx *context, sitePackages FileSystem) vm.ObjectLoader {
	return func(p *vm.Plasma) vm.Value {
		return p.NewFunction(true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					module := arguments[0]
					if _, ok := module.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(module.TypeName(), vm.StringName)
					}
					nameParts := strings.Split(module.GetString(), "@")
					numberOfParts := len(nameParts)
					if numberOfParts > 2 {
						return nil, newModuleNomenclatureError(p, module.GetString())
					}
					moduleName := nameParts[0]
					var version string
					if numberOfParts == 2 {
						version = nameParts[1]
					} else {
						version = "latest"
					}
					// Backup the context
					ctxBackup := *ctx
					oldLocation := sitePackages.RelativePwd()

					sitePackages.ResetPath()

					if !sitePackages.ExistsRelative(moduleName) {
						return nil, newModuleNotFoundError(p, moduleName)
					}
					changeDirectoryError := sitePackages.ChangeDirectoryRelative(moduleName)
					if changeDirectoryError != nil {
						return nil, newChangeDirectoryError(p, changeDirectoryError)
					}
					if version == "latest" {
						moduleVersions, listingError := sitePackages.ListDirectory()
						if listingError != nil {
							return nil, p.NewGoRuntimeError(listingError)
						}
						if len(moduleVersions) == 0 {
							return nil, newNoVersionFoundError(p, moduleName, "latest")
						}
						version = moduleVersions[0]
					}
					changeDirectoryError = sitePackages.ChangeDirectoryRelative(version)
					if changeDirectoryError != nil {
						return nil, newChangeDirectoryError(p, changeDirectoryError)
					}
					// Load the new context
					settingsHandler, openError := sitePackages.OpenRelative("settings.json")
					if openError != nil {
						return nil, p.NewGoRuntimeError(openError)
					}
					jsonContent, readingError := io.ReadAll(settingsHandler)
					if readingError != nil {
						return nil, p.NewGoRuntimeError(readingError)
					}
					var moduleSettings Settings
					jsonParsingError := json.Unmarshal(jsonContent, &moduleSettings)
					if jsonParsingError != nil {
						return nil, p.NewGoRuntimeError(jsonParsingError)
					}
					ctx.moduleName = moduleSettings.Name
					ctx.resources = moduleSettings.Resources
					ctx.version = version
					ctx.entryScript = moduleSettings.EntryScript
					ctx.root = filepath.Join(moduleName, version)
					// Open the entry script
					if !sitePackages.ExistsRelative(ctx.entryScript) {
						return nil, newScriptNotFoundError(p, ctx.entryScript)
					}
					var scriptFile io.ReadSeekCloser
					scriptFile, openError = sitePackages.OpenRelative(ctx.entryScript)
					if openError != nil {
						return nil, p.NewGoRuntimeError(openError)
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return nil, p.NewGoRuntimeError(hashingError)
					}
					if _, ok := memory[scriptHash]; ok {
						return memory[scriptHash], nil
					}
					// Run the entry script
					scriptCode, compilationError := plasma.NewCompiler(reader.NewStringReaderFromFile(scriptFile),
						plasma.Options{
							Debug: false,
						},
					).Compile()
					if compilationError != nil {
						return nil, newCompilationError(p, compilationError)
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(false, p.PeekSymbolTable())
					memory[scriptHash] = script
					p.PushSymbolTable(script.SymbolTable())
					p.PushBytecode(scriptCode)
					_, executionError := p.Execute()
					if executionError != nil {
						return nil, executionError
					}
					p.PopBytecode()
					p.PopSymbolTable()
					// Restore the backed context
					sitePackages.ResetPath()
					sitePackages.ChangeDirectoryFullPath(oldLocation)
					*ctx = ctxBackup
					// Return the module object
					return script, nil
				},
			),
		)
	}
}

func NewImporter(sitePackages FileSystem, pwd FileSystem) map[string]vm.ObjectLoader {
	ctx := &context{
		moduleName:  "",
		version:     "",
		resources:   "",
		entryScript: "",
		root:        "",
	}
	memory := map[string]vm.Value{}
	return map[string]vm.ObjectLoader{
		"open_resource":     getResource(ctx, sitePackages),
		"get_resource_path": getResourcePath(ctx, sitePackages),
		"import_script":     scriptImport(memory, ctx, sitePackages, pwd),
		"import_module":     moduleImport(memory, ctx, sitePackages),

		ResourceReader: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ResourceReader, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 0,
									func(_ vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil
					}))
		},

		ResourceNotFoundError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ResourceNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										resourceName := arguments[0]
										if _, ok := resourceName.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(resourceName.TypeName(), vm.StringName)
										}
										self.SetString(fmt.Sprintf("Resource with name %s not found", resourceName.GetString()))
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		NotInsideModuleError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, NotInsideModuleError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 0,
									func(self vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
										self.SetString("Not inside a module context")
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		ScriptNotFoundError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ScriptNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										script := arguments[0]
										if _, ok := script.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(script.TypeName(), vm.StringName)
										}
										self.SetString(fmt.Sprintf("Script %s not found", script.GetString()))
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil
					},
				),
			)

		},
		CompilationError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, CompilationError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										compilationError := arguments[0]
										if _, ok := compilationError.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(compilationError.TypeName(), vm.StringName)
										}
										self.SetString(compilationError.GetString())
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		ModuleNotFoundError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ModuleNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										moduleName := arguments[0]
										if _, ok := moduleName.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(moduleName.TypeName(), vm.StringName)
										}
										self.SetString(fmt.Sprintf("No module with name %s found", moduleName.GetString()))
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		ChangeDirectoryError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ChangeDirectoryError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										message := arguments[0]
										if _, ok := message.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(message.TypeName(), vm.StringName)
										}
										self.SetString(message.GetString())
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		ModuleNomenclatureError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, ModuleNomenclatureError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 1,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										moduleName := arguments[0]
										if _, ok := moduleName.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(moduleName.TypeName(), vm.StringName)
										}
										self.SetString(fmt.Sprintf("Invalid module nomenclature for %s, expecting estructure like \"NAME\" or \"NAME@VERSION\"", moduleName.GetString()))
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
		NoVersionFoundError: func(p *vm.Plasma) vm.Value {
			return p.NewType(true, NoVersionFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(object vm.Value) *vm.Object {
						object.Set(vm.Initialize,
							p.NewFunction(true, object.SymbolTable(),
								vm.NewBuiltInClassFunction(object, 2,
									func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
										moduleName := arguments[0]
										if _, ok := moduleName.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(moduleName.TypeName(), vm.StringName)
										}
										version := arguments[1]
										if _, ok := version.(*vm.String); !ok {
											return nil, p.NewInvalidTypeError(version.TypeName(), vm.StringName)
										}
										self.SetString(fmt.Sprintf("no module found with name: %s and version %s", moduleName.GetString(), version.GetString()))
										return p.NewNone(), nil
									},
								),
							),
						)
						return nil

					},
				),
			)

		},
	}
}
