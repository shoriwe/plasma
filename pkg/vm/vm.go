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
		true, false, none *Value   // TODO: init me
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
		switch err := recover().(type) {
		case error:
			ctx.err <- err
		case string:
			ctx.err <- fmt.Errorf(err)
		case nil:
			ctx.err <- nil
		default:
			panic("invalid panic, it should be error, string or nil")
		}
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
	ctx.result = make(chan *Value)
	ctx.err = make(chan error)
	ctx.stop = make(chan struct{})
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
