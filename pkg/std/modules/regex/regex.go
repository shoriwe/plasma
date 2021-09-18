package regex

import (
	"github.com/shoriwe/gplasma/pkg/std/features/importlib"
	"github.com/shoriwe/gplasma/pkg/vm"
	"regexp"
)

func prepareRegexp(s *vm.Value, p *vm.Plasma, context *vm.Context) (*regexp.Regexp, *vm.Value) {
	if !s.IsTypeById(vm.StringId) {
		return nil, p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName)
	}
	regex, compilationError := regexp.Compile(s.String)
	if compilationError != nil {
		return nil, p.NewGoRuntimeError(context, compilationError)
	}
	return regex, nil
}

func match(s *vm.Value, regex *regexp.Regexp, p *vm.Plasma, context *vm.Context) (*vm.Value, bool) {
	var target []byte
	if s.IsTypeById(vm.StringId) {
		target = []byte(s.String)
	} else if s.IsTypeById(vm.BytesId) {
		target = s.Bytes
	} else {
		return p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName, vm.BytesName), false
	}
	return p.InterpretAsBool(regex.Match(target)), true
}

func findAll(s *vm.Value, regex *regexp.Regexp, p *vm.Plasma, context *vm.Context) (*vm.Value, bool) {
	isString := false
	if s.IsTypeById(vm.StringId) {
		isString = true
	} else if !s.IsTypeById(vm.BytesId) {
		return p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName, vm.BytesName), false
	}
	var result []*vm.Value
	if isString {
		for _, found := range regex.FindAllString(s.String, -1) {
			result = append(result, p.NewString(context, false, found))
		}
	} else {
		for _, found := range regex.FindAll(s.Bytes, -1) {
			result = append(result, p.NewBytes(context, false, found))
		}
	}
	return p.NewTuple(context, false, result), true
}

func replaceAll(s *vm.Value, rep *vm.Value, regex *regexp.Regexp, p *vm.Plasma, context *vm.Context) (*vm.Value, bool) {
	if s.IsTypeById(vm.StringId) && rep.IsTypeById(vm.StringId) {
		return p.NewString(context, false, regex.ReplaceAllString(s.String, rep.String)), true
	} else if s.IsTypeById(vm.BytesId) && rep.IsTypeById(vm.BytesId) {
		return p.NewBytes(context, false, regex.ReplaceAll(s.Bytes, rep.Bytes)), true
	}
	return p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName, vm.BytesName), false
}

func split(s *vm.Value, regex *regexp.Regexp, p *vm.Plasma, context *vm.Context) (*vm.Value, bool) {
	if !s.IsTypeById(vm.StringId) {
		return p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName, vm.BytesName), false
	}
	var result []*vm.Value
	for _, ss := range regex.Split(s.String, -1) {
		result = append(result, p.NewString(context, false, ss))
	}
	return p.NewTuple(context, false, result), true
}

func submatch(s *vm.Value, regex *regexp.Regexp, p *vm.Plasma, context *vm.Context) (*vm.Value, bool) {
	if s.IsTypeById(vm.StringId) {
		var result []*vm.Value
		for _, ss := range regex.FindStringSubmatch(s.String) {
			result = append(result, p.NewString(context, false, ss))
		}
		return p.NewTuple(context, false, result), true
	} else if s.IsTypeById(vm.BytesId) {
		var result []*vm.Value
		for _, ss := range regex.FindSubmatch(s.Bytes) {
			result = append(result, p.NewBytes(context, false, ss))
		}
		return p.NewTuple(context, false, result), true
	} else {
		return p.NewInvalidTypeError(context, s.GetClass(p).Name, vm.StringName, vm.BytesName), false
	}
}

func RegexpInitialize(p *vm.Plasma) vm.ConstructorCallBack {
	return func(context *vm.Context, value *vm.Value) *vm.Value {
		var (
			regex *regexp.Regexp
		)
		value.Set(p, context, vm.Initialize,
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						var err *vm.Value
						regex, err = prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return p.GetNone(), true
					},
				),
			),
		)
		value.Set(p, context, "Match",
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						return match(arguments[0], regex, p, context)
					},
				),
			),
		)
		value.Set(p, context, "FindAll",
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						return findAll(arguments[0], regex, p, context)
					},
				),
			),
		)
		value.Set(p, context, "ReplaceAll",
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 2,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						return replaceAll(arguments[0], arguments[1], regex, p, context)
					},
				),
			),
		)
		value.Set(p, context, "Split",
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						return split(arguments[0], regex, p, context)
					},
				),
			),
		)
		value.Set(p, context, "FindSubmatch",
			p.NewFunction(context, false, value.SymbolTable(),
				vm.NewBuiltInClassFunction(value, 1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						return submatch(arguments[0], regex, p, context)
					},
				),
			),
		)
		return nil
	}
}

func regexLoader(context *vm.Context, p *vm.Plasma) *vm.Value {
	result := p.NewModule(context, true)
	result.SetOnDemandSymbol("Regexp",
		func() *vm.Value {
			return p.NewType(context, true, "Regexp", result.SymbolTable(), nil,
				vm.NewBuiltInConstructor(RegexpInitialize(p)),
			)
		},
	)
	result.SetOnDemandSymbol("match",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(2,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						regex, err := prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return match(arguments[0], regex, p, context)
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol("findall",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(2,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						regex, err := prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return findAll(arguments[1], regex, p, context)
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol("replaceall",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(3,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						regex, err := prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return replaceAll(arguments[1], arguments[2], regex, p, context)
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol("split",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(2,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						regex, err := prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return split(arguments[1], regex, p, context)
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol("find_submatch",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(2,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						regex, err := prepareRegexp(arguments[0], p, context)
						if err != nil {
							return err, false
						}
						return submatch(arguments[1], regex, p, context)
					},
				),
			)
		},
	)
	return result
}

var Regex = importlib.ModuleInformation{
	Name:   "regexp",
	Loader: regexLoader,
}
