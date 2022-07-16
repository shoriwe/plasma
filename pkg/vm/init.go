package vm

import (
	special_symbols "github.com/shoriwe/gplasma/pkg/common/special-symbols"
)

func (plasma *Plasma) init() {
	// Init classes
	plasma.metaClass()
	plasma.value = plasma.valueClass()
	plasma.string = plasma.stringClass()
	plasma.bytes = plasma.bytesClass()
	plasma.bool = plasma.boolClass()
	plasma.noneType = plasma.noneClass()
	plasma.int = plasma.integerClass()
	plasma.float = plasma.floatClass()
	plasma.array = plasma.arrayClass()
	plasma.tuple = plasma.tupleClass()
	plasma.hash = plasma.hashClass()
	plasma.function = plasma.functionClass()
	// Init values
	plasma.true = plasma.NewBool(true)
	plasma.false = plasma.NewBool(false)
	plasma.none = plasma.NewNone()
	// Init symbols
	// -- Classes
	plasma.rootSymbols.Set(special_symbols.Value, plasma.value)
	plasma.rootSymbols.Set(special_symbols.String, plasma.string)
	plasma.rootSymbols.Set(special_symbols.Bytes, plasma.bytes)
	plasma.rootSymbols.Set(special_symbols.Bool, plasma.bool)
	plasma.rootSymbols.Set(special_symbols.None, plasma.noneType)
	plasma.rootSymbols.Set(special_symbols.Int, plasma.int)
	plasma.rootSymbols.Set(special_symbols.Float, plasma.float)
	plasma.rootSymbols.Set(special_symbols.Array, plasma.array)
	plasma.rootSymbols.Set(special_symbols.Tuple, plasma.tuple)
	plasma.rootSymbols.Set(special_symbols.Hash, plasma.hash)
	plasma.rootSymbols.Set(special_symbols.Function, plasma.function)
	plasma.rootSymbols.Set(special_symbols.Class, plasma.class)
	/*
		- print
		- println
		- range TODO
	*/
	plasma.rootSymbols.Set(special_symbols.Print, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			for index, arg := range argument {
				if index != 0 {
					_, writeError := plasma.Stdout.Write([]byte(" "))
					if writeError != nil {
						panic(writeError)
					}
				}
				_, writeError := plasma.Stdout.Write([]byte(arg.String()))
				if writeError != nil {
					panic(writeError)
				}

			}
			return plasma.none, nil
		},
	))
	plasma.rootSymbols.Set(special_symbols.Println, plasma.NewBuiltInFunction(plasma.rootSymbols,
		func(argument ...*Value) (*Value, error) {
			for index, arg := range argument {
				if index != 0 {
					_, writeError := plasma.Stdout.Write([]byte(" "))
					if writeError != nil {
						panic(writeError)
					}
				}
				_, writeError := plasma.Stdout.Write([]byte(arg.String()))
				if writeError != nil {
					panic(writeError)
				}

			}
			_, writeError := plasma.Stdout.Write([]byte("\n"))
			if writeError != nil {
				panic(writeError)
			}
			return plasma.none, nil
		},
	))
}
