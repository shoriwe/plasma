package vm

import (
	"bufio"
	"crypto/rand"
	"hash"
	"hash/crc32"
	"io"
	"math/big"
)

const (
	polySize = 0xffffffff
)

type ObjectLoader func(*Context, *Plasma) Value

type Plasma struct {
	currentId      int64
	builtInContext *Context
	Crc32Hash      hash.Hash32
	seed           uint64
	stdinScanner   *bufio.Scanner
	stdin          io.Reader
	stdout         io.Writer
	stderr         io.Writer
}

func (p *Plasma) HashString(s string) int64 {
	_, hashingError := p.Crc32Hash.Write([]byte(s))
	if hashingError != nil {
		panic(hashingError)
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue)
}

func (p *Plasma) HashBytes(s []byte) int64 {
	_, hashingError := p.Crc32Hash.Write(s)
	if hashingError != nil {
		panic(hashingError)
	}
	hashValue := p.Crc32Hash.Sum32()
	p.Crc32Hash.Reset()
	return int64(hashValue)
}

func (p *Plasma) Seed() uint64 {
	return p.seed
}

/*
	LoadBuiltInObject
	This function should be used to load custom object in the built-in symbol table
*/
func (p *Plasma) LoadBuiltInObject(symbolName string, loader ObjectLoader) {
	p.builtInContext.PeekSymbolTable().Set(symbolName, loader(p.builtInContext, p))
}

func (p *Plasma) LoadBuiltInSymbols(symbolMap map[string]ObjectLoader) {
	for symbol, loader := range symbolMap {
		p.builtInContext.PeekSymbolTable().Set(symbol, loader(p.builtInContext, p))
	}
}

/*
	InitializeBytecode
	Loads the bytecode and clears the stack
*/

func (p *Plasma) StdInScanner() *bufio.Scanner {
	return p.stdinScanner
}

func (p *Plasma) StdIn() io.Reader {
	return p.stdin
}

func (p *Plasma) StdOut() io.Writer {
	return p.stdout
}

func (p *Plasma) StdErr() io.Writer {
	return p.stderr
}

func (p *Plasma) BuiltInSymbols() *SymbolTable {
	return p.builtInContext.PeekSymbolTable()
}

func (p *Plasma) NextId() int64 {
	result := p.currentId
	p.currentId++
	return result
}

func (p *Plasma) InitializeContext(context *Context) {
	symbols := NewSymbolTable(p.builtInContext.PeekSymbolTable())
	symbols.Set("__built_in__",
		&Object{
			id:         p.NextId(),
			typeName:   ObjectName,
			class:      nil,
			subClasses: nil,
			isBuiltIn:  true,
			symbols:    p.builtInContext.PeekSymbolTable(),
		},
	)
	symbols.Set(Self,
		&Object{
			id:         p.NextId(),
			typeName:   ObjectName,
			class:      nil,
			subClasses: nil,
			symbols:    symbols,
		},
	)
	context.SymbolStack.Push(symbols)
}

func NewPlasmaVM(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Plasma {
	number, randError := rand.Int(rand.Reader, big.NewInt(polySize))
	if randError != nil {
		panic(randError)
	}
	vm := &Plasma{
		currentId:      1,
		builtInContext: NewContext(),
		Crc32Hash:      crc32.New(crc32.MakeTable(polySize)),
		seed:           number.Uint64(),
		stdinScanner:   bufio.NewScanner(stdin),
		stdin:          stdin,
		stdout:         stdout,
		stderr:         stderr,
	}
	vm.builtInContext.PushSymbolTable(NewSymbolTable(nil))
	vm.InitializeBuiltIn()
	return vm
}
