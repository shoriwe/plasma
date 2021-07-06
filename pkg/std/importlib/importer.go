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

type importContext struct {
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
	return func(context *vm.Context, object vm.Value) *vm.Object {
		object.SetOnDemandSymbol("Read",
			func() vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 1,
						func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
							bytesToRead := arguments[0]
							if _, ok := bytesToRead.(*vm.Integer); !ok {
								return p.NewInvalidTypeError(context, bytesToRead.TypeName(), vm.IntegerName), nil
							}
							bytes := make([]byte, bytesToRead.GetInteger())
							numberOfBytes, readError := r.Read(bytes)
							if readError != nil {
								if readError == io.EOF {
									return p.GetNone(), nil
								} else {
									return nil, p.NewGoRuntimeError(context, readError)
								}
							}
							bytes = bytes[:numberOfBytes]
							return p.NewBytes(context, false, context.PeekSymbolTable(), bytes), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol("Seek",
			func() vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 1,
						func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
							seek := arguments[0]
							if _, ok := seek.(*vm.Integer); !ok {
								return p.NewInvalidTypeError(context, seek.TypeName(), vm.IntegerName), nil
							}
							_, seekError := r.Seek(seek.GetInteger(), io.SeekStart)
							if seekError != nil {
								return nil, p.NewGoRuntimeError(context, seekError)
							}
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol("Close",
			func() vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 0,
						func(self vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
							closeError := r.Close()
							if closeError != nil {
								return nil, p.NewGoRuntimeError(context, closeError)
							}
							return p.GetNone(), nil
						},
					),
				)
			},
		)
		return nil
	}
}

func newResourceReader(context *vm.Context, p *vm.Plasma, r io.ReadSeekCloser) vm.Value {
	resourceReader := p.NewObject(context, false, ResourceReader, nil,
		context.PeekSymbolTable(),
	)
	resourceReaderInitialize(p, r)(context, resourceReader)
	return resourceReader
}

func newResourceNotFoundError(context *vm.Context, p *vm.Plasma, path string) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ResourceNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), path),
	)
	return result.(*vm.Object)
}

func newNotInsideModuleError(context *vm.Context, p *vm.Plasma) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(NotInsideModuleError))
	p.ForceInitialization(context, result)
	return result.(*vm.Object)
}

func newScriptNotFoundError(context *vm.Context, p *vm.Plasma, path string) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ScriptNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), path),
	)
	return result.(*vm.Object)
}

func newCompilationError(context *vm.Context, p *vm.Plasma, compilationError *errors.Error) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(CompilationError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), compilationError.Message()),
	)
	return result.(*vm.Object)
}

func newModuleNotFoundError(context *vm.Context, p *vm.Plasma, moduleName string) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ModuleNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
	)
	return result.(*vm.Object)
}

func newChangeDirectoryError(context *vm.Context, p *vm.Plasma, compilationError *errors.Error) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ChangeDirectoryError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), compilationError.Message()),
	)
	return result.(*vm.Object)
}

func newModuleNomenclatureError(context *vm.Context, p *vm.Plasma, moduleName string) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ModuleNomenclatureError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
	)
	return result.(*vm.Object)
}

func newNoVersionFoundError(context *vm.Context, p *vm.Plasma, moduleName string, version string) *vm.Object {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(NoVersionFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
		p.NewString(context, false, context.PeekSymbolTable(), version),
	)
	return result.(*vm.Object)
}

func (c *importContext) isSet() bool {
	return c.root != ""
}

func getResource(ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(_ vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					if !ctx.isSet() {
						return nil, newNotInsideModuleError(context, p)
					}

					resourcePathObject := arguments[0]
					if _, ok := resourcePathObject.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(context, resourcePathObject.TypeName(), vm.StringName)
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					resourcePath := resourcePathObject.GetString()
					if !sitePackages.ExistsRelative(resourcePath) {
						return nil, newResourceNotFoundError(context, p, resourcePath)
					}
					resourceHandler, openError := sitePackages.OpenRelative(resourcePath)
					if openError != nil {
						return nil, p.NewGoRuntimeError(context, openError)
					}
					return newResourceReader(context, p, resourceHandler), nil
				},
			),
		)
	}
}

func getResourcePath(ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					resource := arguments[0]
					if _, ok := resource.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(context, resource.TypeName(), vm.StringName)
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					if !sitePackages.ExistsRelative(resource.GetString()) {
						return nil, newResourceNotFoundError(context, p, resource.GetString())
					}
					resourcePath := filepath.Join(sitePackages.AbsolutePwd(), ctx.root, ctx.resources, resource.GetString())
					return p.NewString(context, false, context.PeekSymbolTable(), resourcePath), nil
				},
			),
		)
	}
}

func scriptImport(memory map[string]vm.Value, ctx *importContext, sitePackages FileSystem, pwd FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					scriptPath := arguments[0]
					if _, ok := scriptPath.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(context, scriptPath.TypeName(), vm.StringName)
					}
					// First try to get the content of the file
					var (
						scriptFile io.ReadSeekCloser
						openError  error
					)
					// Check if the code is being ran in the importContext of a module
					if ctx.isSet() {
						// If it is, import the scriptFile relative the root of the module
						// Check if file exists
						if !sitePackages.ExistsRelative(scriptPath.GetString()) {
							return nil, newScriptNotFoundError(context, p, scriptPath.GetString())
						}
						scriptFile, openError = sitePackages.OpenRelative(scriptPath.GetString())

						oldLocation := sitePackages.RelativePwd()
						defer sitePackages.ChangeDirectoryFullPath(oldLocation)
						sitePackages.ChangeDirectoryToFileLocation(scriptPath.GetString())
					} else {
						// If not import the scriptFile from the immediate filesystem of the running scriptFile
						// Check if file exists
						if !pwd.ExistsRelative(scriptPath.GetString()) {
							return nil, newScriptNotFoundError(context, p, scriptPath.GetString())
						}
						scriptFile, openError = pwd.OpenRelative(scriptPath.GetString())
						oldLocation := pwd.RelativePwd()
						defer pwd.ChangeDirectoryFullPath(oldLocation)
						pwd.ChangeDirectoryToFileLocation(scriptPath.GetString())
					}
					if openError != nil {
						return nil, p.NewGoRuntimeError(context, openError)
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return nil, p.NewGoRuntimeError(context, hashingError)
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
						return nil, newCompilationError(context, p, compilationError)
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(context, false, context.PeekSymbolTable())
					memory[scriptHash] = script
					context.PushSymbolTable(script.SymbolTable())
					_, executionError := p.Execute(context, scriptCode)
					if executionError != nil {
						return nil, executionError
					}
					context.PopSymbolTable()
					// Return the initialized module object
					return script, nil
				},
			),
		)
	}
}

func moduleImport(memory map[string]vm.Value, ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
					module := arguments[0]
					if _, ok := module.(*vm.String); !ok {
						return nil, p.NewInvalidTypeError(context, module.TypeName(), vm.StringName)
					}
					nameParts := strings.Split(module.GetString(), "@")
					numberOfParts := len(nameParts)
					if numberOfParts > 2 {
						return nil, newModuleNomenclatureError(context, p, module.GetString())
					}
					moduleName := nameParts[0]
					var version string
					if numberOfParts == 2 {
						version = nameParts[1]
					} else {
						version = "latest"
					}
					// Backup the importContext
					ctxBackup := *ctx
					oldLocation := sitePackages.RelativePwd()

					sitePackages.ResetPath()

					if !sitePackages.ExistsRelative(moduleName) {
						return nil, newModuleNotFoundError(context, p, moduleName)
					}
					changeDirectoryError := sitePackages.ChangeDirectoryRelative(moduleName)
					if changeDirectoryError != nil {
						return nil, newChangeDirectoryError(context, p, changeDirectoryError)
					}
					if version == "latest" {
						moduleVersions, listingError := sitePackages.ListDirectory()
						if listingError != nil {
							return nil, p.NewGoRuntimeError(context, listingError)
						}
						if len(moduleVersions) == 0 {
							return nil, newNoVersionFoundError(context, p, moduleName, "latest")
						}
						version = moduleVersions[0]
					}
					changeDirectoryError = sitePackages.ChangeDirectoryRelative(version)
					if changeDirectoryError != nil {
						return nil, newChangeDirectoryError(context, p, changeDirectoryError)
					}
					// Load the new importContext
					settingsHandler, openError := sitePackages.OpenRelative("settings.json")
					if openError != nil {
						return nil, p.NewGoRuntimeError(context, openError)
					}
					jsonContent, readingError := io.ReadAll(settingsHandler)
					if readingError != nil {
						return nil, p.NewGoRuntimeError(context, readingError)
					}
					var moduleSettings Settings
					jsonParsingError := json.Unmarshal(jsonContent, &moduleSettings)
					if jsonParsingError != nil {
						return nil, p.NewGoRuntimeError(context, jsonParsingError)
					}
					ctx.moduleName = moduleSettings.Name
					ctx.resources = moduleSettings.Resources
					ctx.version = version
					ctx.entryScript = moduleSettings.EntryScript
					ctx.root = filepath.Join(moduleName, version)
					// Open the entry script
					if !sitePackages.ExistsRelative(ctx.entryScript) {
						return nil, newScriptNotFoundError(context, p, ctx.entryScript)
					}
					var scriptFile io.ReadSeekCloser
					scriptFile, openError = sitePackages.OpenRelative(ctx.entryScript)
					if openError != nil {
						return nil, p.NewGoRuntimeError(context, openError)
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return nil, p.NewGoRuntimeError(context, hashingError)
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
						return nil, newCompilationError(context, p, compilationError)
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(context, false, context.PeekSymbolTable())
					memory[scriptHash] = script
					context.PushSymbolTable(script.SymbolTable())
					_, executionError := p.Execute(context, scriptCode)
					if executionError != nil {
						return nil, executionError
					}
					context.PopSymbolTable()
					// Restore the backed importContext
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
	ctx := &importContext{
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

		ResourceReader: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ResourceReader, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 0,
										func(_ vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil
					},
				),
			)
		},

		ResourceNotFoundError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ResourceNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											resourceName := arguments[0]
											if _, ok := resourceName.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, resourceName.TypeName(), vm.StringName)
											}
											self.SetString(fmt.Sprintf("Resource with name %s not found", resourceName.GetString()))
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		NotInsideModuleError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, NotInsideModuleError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 0,
										func(self vm.Value, _ ...vm.Value) (vm.Value, *vm.Object) {
											self.SetString("Not inside a module importContext")
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		ScriptNotFoundError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ScriptNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											script := arguments[0]
											if _, ok := script.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, script.TypeName(), vm.StringName)
											}
											self.SetString(fmt.Sprintf("Script %s not found", script.GetString()))
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil
					},
				),
			)

		},
		CompilationError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, CompilationError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											compilationError := arguments[0]
											if _, ok := compilationError.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, compilationError.TypeName(), vm.StringName)
											}
											self.SetString(compilationError.GetString())
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		ModuleNotFoundError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ModuleNotFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											moduleName := arguments[0]
											if _, ok := moduleName.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName)
											}
											self.SetString(fmt.Sprintf("No module with name %s found", moduleName.GetString()))
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		ChangeDirectoryError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ChangeDirectoryError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											message := arguments[0]
											if _, ok := message.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, message.TypeName(), vm.StringName)
											}
											self.SetString(message.GetString())
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		ModuleNomenclatureError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, ModuleNomenclatureError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											moduleName := arguments[0]
											if _, ok := moduleName.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName)
											}
											self.SetString(fmt.Sprintf("Invalid module nomenclature for %s, expecting estructure like \"NAME\" or \"NAME@VERSION\"", moduleName.GetString()))
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
		NoVersionFoundError: func(context *vm.Context, p *vm.Plasma) vm.Value {
			return p.NewType(context, true, NoVersionFoundError, p.BuiltInSymbols(), []*vm.Type{p.ForceMasterGetAny(vm.TypeName).(*vm.Type)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object vm.Value) *vm.Object {
						object.SetOnDemandSymbol(vm.Initialize,
							func() vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 2,
										func(self vm.Value, arguments ...vm.Value) (vm.Value, *vm.Object) {
											moduleName := arguments[0]
											if _, ok := moduleName.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName)
											}
											version := arguments[1]
											if _, ok := version.(*vm.String); !ok {
												return nil, p.NewInvalidTypeError(context, version.TypeName(), vm.StringName)
											}
											self.SetString(fmt.Sprintf("no module found with name: %s and version %s", moduleName.GetString(), version.GetString()))
											return p.GetNone(), nil
										},
									),
								)
							},
						)
						return nil

					},
				),
			)

		},
	}
}
