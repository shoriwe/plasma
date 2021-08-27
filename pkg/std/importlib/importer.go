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
	return func(context *vm.Context, object *vm.Value) *vm.Value {
		object.SetOnDemandSymbol("Read",
			func() *vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 1,
						func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
							bytesToRead := arguments[0]
							if !bytesToRead.IsTypeById(vm.IntegerId) {
								return p.NewInvalidTypeError(context, bytesToRead.TypeName(), vm.IntegerName), false
							}
							bytes := make([]byte, bytesToRead.Integer)
							numberOfBytes, readError := r.Read(bytes)
							if readError != nil {
								if readError == io.EOF {
									return p.GetNone(), true
								} else {
									return p.NewGoRuntimeError(context, readError), false
								}
							}
							bytes = bytes[:numberOfBytes]
							return p.NewBytes(context, false, context.PeekSymbolTable(), bytes), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol("Seek",
			func() *vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 1,
						func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
							seek := arguments[0]
							if !seek.IsTypeById(vm.IntegerId) {
								return p.NewInvalidTypeError(context, seek.TypeName(), vm.IntegerName), false
							}
							_, seekError := r.Seek(seek.Integer, io.SeekStart)
							if seekError != nil {
								return p.NewGoRuntimeError(context, seekError), false
							}
							return p.GetNone(), true
						},
					),
				)
			},
		)
		object.SetOnDemandSymbol("Close",
			func() *vm.Value {
				return p.NewFunction(context, false, object.SymbolTable(),
					vm.NewBuiltInClassFunction(object, 0,
						func(self *vm.Value, _ ...*vm.Value) (*vm.Value, bool) {
							closeError := r.Close()
							if closeError != nil {
								return p.NewGoRuntimeError(context, closeError), false
							}
							return p.GetNone(), true
						},
					),
				)
			},
		)
		return nil
	}
}

func newResourceReader(context *vm.Context, p *vm.Plasma, r io.ReadSeekCloser) *vm.Value {
	resourceReader := p.NewValue(context, false, ResourceReader, nil,
		context.PeekSymbolTable(),
	)
	resourceReaderInitialize(p, r)(context, resourceReader)
	return resourceReader
}

func newResourceNotFoundError(context *vm.Context, p *vm.Plasma, path string) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ResourceNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), path),
	)
	return result
}

func newNotInsideModuleError(context *vm.Context, p *vm.Plasma) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(NotInsideModuleError))
	p.ForceInitialization(context, result)
	return result
}

func newScriptNotFoundError(context *vm.Context, p *vm.Plasma, path string) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ScriptNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), path),
	)
	return result
}

func newCompilationError(context *vm.Context, p *vm.Plasma, compilationError *errors.Error) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(CompilationError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), compilationError.Message()),
	)
	return result
}

func newModuleNotFoundError(context *vm.Context, p *vm.Plasma, moduleName string) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ModuleNotFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
	)
	return result
}

func newChangeDirectoryError(context *vm.Context, p *vm.Plasma, compilationError *errors.Error) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ChangeDirectoryError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), compilationError.Message()),
	)
	return result
}

func newModuleNomenclatureError(context *vm.Context, p *vm.Plasma, moduleName string) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(ModuleNomenclatureError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
	)
	return result
}

func newNoVersionFoundError(context *vm.Context, p *vm.Plasma, moduleName string, version string) *vm.Value {
	result := p.ForceConstruction(context, p.ForceMasterGetAny(NoVersionFoundError))
	p.ForceInitialization(context, result,
		p.NewString(context, false, context.PeekSymbolTable(), moduleName),
		p.NewString(context, false, context.PeekSymbolTable(), version),
	)
	return result
}

func (c *importContext) isSet() bool {
	return c.root != ""
}

func getResource(ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) *vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(_ *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
					if !ctx.isSet() {
						return newNotInsideModuleError(context, p), false
					}

					resourcePathObject := arguments[0]
					if !resourcePathObject.IsTypeById(vm.StringId) {
						return p.NewInvalidTypeError(context, resourcePathObject.TypeName(), vm.StringName), false
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					resourcePath := resourcePathObject.String
					if !sitePackages.ExistsRelative(resourcePath) {
						return newResourceNotFoundError(context, p, resourcePath), false
					}
					resourceHandler, openError := sitePackages.OpenRelative(resourcePath)
					if openError != nil {
						return p.NewGoRuntimeError(context, openError), false
					}
					return newResourceReader(context, p, resourceHandler), true
				},
			),
		)
	}
}

func getResourcePath(ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) *vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
					resource := arguments[0]
					if !resource.IsTypeById(vm.StringId) {
						return p.NewInvalidTypeError(context, resource.TypeName(), vm.StringName), false
					}

					oldLocation := sitePackages.RelativePwd()
					defer sitePackages.ChangeDirectoryFullPath(oldLocation)
					sitePackages.ChangeDirectoryRelative(ctx.resources)

					if !sitePackages.ExistsRelative(resource.String) {
						return newResourceNotFoundError(context, p, resource.String), false
					}
					resourcePath := filepath.Join(sitePackages.AbsolutePwd(), ctx.root, ctx.resources, resource.String)
					return p.NewString(context, false, context.PeekSymbolTable(), resourcePath), true
				},
			),
		)
	}
}

func scriptImport(memory map[string]*vm.Value, ctx *importContext, sitePackages FileSystem, pwd FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) *vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
					scriptPath := arguments[0]
					if !scriptPath.IsTypeById(vm.StringId) {
						return p.NewInvalidTypeError(context, scriptPath.TypeName(), vm.StringName), false
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
						if !sitePackages.ExistsRelative(scriptPath.String) {
							return newScriptNotFoundError(context, p, scriptPath.String), false
						}
						scriptFile, openError = sitePackages.OpenRelative(scriptPath.String)

						oldLocation := sitePackages.RelativePwd()
						defer sitePackages.ChangeDirectoryFullPath(oldLocation)
						sitePackages.ChangeDirectoryToFileLocation(scriptPath.String)
					} else {
						// If not import the scriptFile from the immediate filesystem of the running scriptFile
						// Check if file exists
						if !pwd.ExistsRelative(scriptPath.String) {
							return newScriptNotFoundError(context, p, scriptPath.String), false
						}
						scriptFile, openError = pwd.OpenRelative(scriptPath.String)
						oldLocation := pwd.RelativePwd()
						defer pwd.ChangeDirectoryFullPath(oldLocation)
						pwd.ChangeDirectoryToFileLocation(scriptPath.String)
					}
					if openError != nil {
						return p.NewGoRuntimeError(context, openError), false
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return p.NewGoRuntimeError(context, hashingError), false
					}
					if _, ok := memory[scriptHash]; ok {
						return memory[scriptHash], true
					}
					// ToDo: Fix this, use a better file reader object
					scriptCode, compilationError := plasma.NewCompiler(reader.NewStringReaderFromFile(scriptFile),
						plasma.Options{
							Debug: false,
						},
					).Compile()
					if compilationError != nil {
						return newCompilationError(context, p, compilationError), false
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(context, false, context.PeekSymbolTable())
					memory[scriptHash] = script
					context.PushSymbolTable(script.SymbolTable())
					executionError, success := p.Execute(context, scriptCode)
					if !success {
						return executionError, false
					}
					context.PopSymbolTable()
					// Return the initialized module object
					return script, true
				},
			),
		)
	}
}

func moduleImport(memory map[string]*vm.Value, ctx *importContext, sitePackages FileSystem) vm.ObjectLoader {
	return func(context *vm.Context, p *vm.Plasma) *vm.Value {
		return p.NewFunction(context, true, p.BuiltInSymbols(),
			vm.NewBuiltInFunction(1,
				func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
					module := arguments[0]
					if !module.IsTypeById(vm.StringId) {
						return p.NewInvalidTypeError(context, module.TypeName(), vm.StringName), false
					}
					nameParts := strings.Split(module.String, "@")
					numberOfParts := len(nameParts)
					if numberOfParts > 2 {
						return newModuleNomenclatureError(context, p, module.String), false
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
						return newModuleNotFoundError(context, p, moduleName), false
					}
					changeDirectoryError := sitePackages.ChangeDirectoryRelative(moduleName)
					if changeDirectoryError != nil {
						return newChangeDirectoryError(context, p, changeDirectoryError), false
					}
					if version == "latest" {
						moduleVersions, listingError := sitePackages.ListDirectory()
						if listingError != nil {
							return p.NewGoRuntimeError(context, listingError), false
						}
						if len(moduleVersions) == 0 {
							return newNoVersionFoundError(context, p, moduleName, "latest"), false
						}
						version = moduleVersions[0]
					}
					changeDirectoryError = sitePackages.ChangeDirectoryRelative(version)
					if changeDirectoryError != nil {
						return newChangeDirectoryError(context, p, changeDirectoryError), false
					}
					// Load the new importContext
					settingsHandler, openError := sitePackages.OpenRelative("settings.json")
					if openError != nil {
						return p.NewGoRuntimeError(context, openError), false
					}
					jsonContent, readingError := io.ReadAll(settingsHandler)
					if readingError != nil {
						return p.NewGoRuntimeError(context, readingError), false
					}
					var moduleSettings Settings
					jsonParsingError := json.Unmarshal(jsonContent, &moduleSettings)
					if jsonParsingError != nil {
						return p.NewGoRuntimeError(context, jsonParsingError), false
					}
					ctx.moduleName = moduleSettings.Name
					ctx.resources = moduleSettings.Resources
					ctx.version = version
					ctx.entryScript = moduleSettings.EntryScript
					ctx.root = filepath.Join(moduleName, version)
					// Open the entry script
					if !sitePackages.ExistsRelative(ctx.entryScript) {
						return newScriptNotFoundError(context, p, ctx.entryScript), false
					}
					var scriptFile io.ReadSeekCloser
					scriptFile, openError = sitePackages.OpenRelative(ctx.entryScript)
					if openError != nil {
						return p.NewGoRuntimeError(context, openError), false
					}
					// Check if the script was already imported
					scriptHash, hashingError := getScriptHash(scriptFile)
					if hashingError != nil {
						return p.NewGoRuntimeError(context, hashingError), false
					}
					if _, ok := memory[scriptHash]; ok {
						return memory[scriptHash], true
					}
					// Run the entry script
					scriptCode, compilationError := plasma.NewCompiler(reader.NewStringReaderFromFile(scriptFile),
						plasma.Options{
							Debug: false,
						},
					).Compile()
					if compilationError != nil {
						return newCompilationError(context, p, compilationError), false
					}
					// Prepare the module object that will receive the namespace
					script := p.NewModule(context, false, context.PeekSymbolTable())
					memory[scriptHash] = script
					context.PushSymbolTable(script.SymbolTable())
					executionError, success := p.Execute(context, scriptCode)
					if !success {
						return executionError, false
					}
					context.PopSymbolTable()
					// Restore the backed importContext
					sitePackages.ResetPath()
					sitePackages.ChangeDirectoryFullPath(oldLocation)
					*ctx = ctxBackup
					// Return the module object
					return script, true
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
	memory := map[string]*vm.Value{}
	return map[string]vm.ObjectLoader{
		"open_resource":     getResource(ctx, sitePackages),
		"get_resource_path": getResourcePath(ctx, sitePackages),
		"import_script":     scriptImport(memory, ctx, sitePackages, pwd),
		"import_module":     moduleImport(memory, ctx, sitePackages),

		ResourceReader: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ResourceReader, p.BuiltInSymbols(), nil,
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 0,
										func(_ *vm.Value, _ ...*vm.Value) (*vm.Value, bool) {
											return p.GetNone(), true
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

		ResourceNotFoundError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ResourceNotFoundError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											resourceName := arguments[0]
											if !resourceName.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, resourceName.TypeName(), vm.StringName), false
											}
											self.SetString(fmt.Sprintf("Resource with name %s not found", resourceName.String))
											return p.GetNone(), true
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
		NotInsideModuleError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, NotInsideModuleError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 0,
										func(self *vm.Value, _ ...*vm.Value) (*vm.Value, bool) {
											self.SetString("Not inside a module importContext")
											return p.GetNone(), true
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
		ScriptNotFoundError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ScriptNotFoundError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											script := arguments[0]
											if !script.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, script.TypeName(), vm.StringName), false
											}
											self.SetString(fmt.Sprintf("Script %s not found", script.String))
											return p.GetNone(), true
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
		CompilationError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, CompilationError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											compilationError := arguments[0]
											if !compilationError.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, compilationError.TypeName(), vm.StringName), false
											}
											self.SetString(compilationError.String)
											return p.GetNone(), true
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
		ModuleNotFoundError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ModuleNotFoundError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											moduleName := arguments[0]
											if !moduleName.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName), false
											}
											self.SetString(fmt.Sprintf("No module with name %s found", moduleName.String))
											return p.GetNone(), true
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
		ChangeDirectoryError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ChangeDirectoryError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											message := arguments[0]
											if !message.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, message.TypeName(), vm.StringName), false
											}
											self.SetString(message.String)
											return p.GetNone(), true
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
		ModuleNomenclatureError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, ModuleNomenclatureError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 1,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											moduleName := arguments[0]
											if !moduleName.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName), false
											}
											self.SetString(fmt.Sprintf("Invalid module nomenclature for %s, expecting estructure like \"NAME\" or \"NAME@VERSION\"", moduleName.String))
											return p.GetNone(), true
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
		NoVersionFoundError: func(context *vm.Context, p *vm.Plasma) *vm.Value {
			return p.NewType(context, true, NoVersionFoundError, p.BuiltInSymbols(), []*vm.Value{p.ForceMasterGetAny(vm.RuntimeError)},
				vm.NewBuiltInConstructor(
					func(context *vm.Context, object *vm.Value) *vm.Value {
						object.SetOnDemandSymbol(vm.Initialize,
							func() *vm.Value {
								return p.NewFunction(context, true, object.SymbolTable(),
									vm.NewBuiltInClassFunction(object, 2,
										func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
											moduleName := arguments[0]
											if !moduleName.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, moduleName.TypeName(), vm.StringName), false
											}
											version := arguments[1]
											if !version.IsTypeById(vm.StringId) {
												return p.NewInvalidTypeError(context, version.TypeName(), vm.StringName), false
											}
											self.SetString(fmt.Sprintf("no module found with name: %s and version %s", moduleName.String, version.String))
											return p.GetNone(), true
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
