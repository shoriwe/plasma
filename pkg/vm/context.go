package vm

import (
	"github.com/shoriwe/gplasma/pkg/common"
)

type (
	contextCode struct {
		bytecode []byte
		index    int64
		onExit   *common.ListStack[[]byte]
	}
	context struct {
		result         chan *Value
		err            chan error
		stop           chan struct{}
		code           *common.ListStack[*contextCode]
		stack          *common.ListStack[*Value]
		register       *Value
		currentSymbols *Symbols
	}
)

func (ctx *context) hasNext() bool {
	for ctx.code.HasNext() {
		ctxCode := ctx.code.Peek()
		if ctxCode.index >= int64(len(ctxCode.bytecode)) {
			ctx.popBlock()
			continue
		}
		return true
	}
	return false
}

func (plasma *Plasma) newContext(bytecode []byte) *context {
	codeStack := &common.ListStack[*contextCode]{}
	codeStack.Push(&contextCode{
		bytecode: bytecode,
		index:    0,
		onExit:   &common.ListStack[[]byte]{},
	})
	return &context{
		result:         nil,
		err:            nil,
		stop:           nil,
		code:           codeStack,
		stack:          &common.ListStack[*Value]{},
		register:       nil,
		currentSymbols: NewSymbols(plasma.rootSymbols),
	}
}
