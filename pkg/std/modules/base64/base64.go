package base64

import (
	"encoding/base64"
	"github.com/shoriwe/gplasma/pkg/std/features/importlib"
	"github.com/shoriwe/gplasma/pkg/vm"
)

func base64Loader(context *vm.Context, p *vm.Plasma) *vm.Value {
	result := p.NewModule(context, true)
	result.SetOnDemandSymbol("encode",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						var src []byte
						if arguments[0].IsTypeById(vm.StringId) {
							src = []byte(arguments[0].String)
						} else if arguments[0].IsTypeById(vm.BytesId) {
							src = arguments[0].Bytes
						} else {
							return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, vm.StringName, vm.BytesName), false
						}
						encoded := make([]byte, base64.StdEncoding.EncodedLen(len(src)))
						base64.StdEncoding.Encode(encoded, src)
						return p.NewBytes(context, false, encoded), true
					},
				),
			)
		},
	)
	result.SetOnDemandSymbol("decode",
		func() *vm.Value {
			return p.NewFunction(context, true, result.SymbolTable(),
				vm.NewBuiltInFunction(1,
					func(self *vm.Value, arguments ...*vm.Value) (*vm.Value, bool) {
						var src []byte
						if arguments[0].IsTypeById(vm.StringId) {
							src = []byte(arguments[0].String)
						} else if arguments[0].IsTypeById(vm.BytesId) {
							src = arguments[0].Bytes
						} else {
							return p.NewInvalidTypeError(context, arguments[0].GetClass(p).Name, vm.StringName, vm.BytesName), false
						}
						decoded := make([]byte, base64.StdEncoding.DecodedLen(len(src)))
						_, decodeError := base64.StdEncoding.Decode(decoded, src)
						if decodeError != nil {
							return p.NewGoRuntimeError(context, decodeError), false
						}
						return p.NewBytes(context, false, decoded), true
					},
				),
			)
		},
	)
	return result
}

var Base64 = importlib.ModuleInformation{
	Name:   "base64",
	Loader: base64Loader,
}
