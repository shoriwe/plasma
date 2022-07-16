package vm

import (
	"fmt"
	"io"
)

type (
	Plasma struct {
		Stdin             io.Reader
		Stdout, Stderr    io.Writer
		rootSymbols       *Symbols // TODO: init me
		onDemand          map[string]func(self *Value) *Value
		true, false, none *Value
		value             *Value
		string            *Value
		bytes             *Value
		bool              *Value
		noneType          *Value
		int               *Value
		float             *Value
		array             *Value
		tuple             *Value
		hash              *Value
		function          *Value
		class             *Value
	}
)

func (plasma *Plasma) executeCtx(ctx *context) {
	defer func() {
		err := recover()
		if err != nil {
			ctx.err <- fmt.Errorf("execution error: %v", err)
		}
		return
	}()
	for ctx.hasNext() {
		select {
		case <-ctx.stop:
			return
		default:
			plasma.do(ctx)
		}
	}
}

func (plasma *Plasma) Execute(bytecode []byte) (result chan *Value, err chan error, stop chan struct{}) {
	// Create new context
	ctx := plasma.newContext(bytecode)
	ctx.result = make(chan *Value, 1)
	ctx.err = make(chan error, 1)
	ctx.stop = make(chan struct{}, 1)
	// Execute bytecode with context
	plasma.executeCtx(ctx)
	return ctx.result, ctx.err, ctx.stop
}

func NewVM(stdin io.Reader, stdout, stderr io.Writer) *Plasma {
	plasma := &Plasma{
		Stdin:       stdin,
		Stdout:      stdout,
		Stderr:      stderr,
		rootSymbols: NewSymbols(nil),
	}
	plasma.init()
	return plasma
}
