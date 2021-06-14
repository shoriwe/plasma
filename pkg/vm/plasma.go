package vm

import (
	"bufio"
	"crypto/rand"
	"hash"
	"hash/crc32"
	"io"
	"math/big"
	"sync"
)

const (
	polySize = 0xffffffff
)

type ObjectLoader func(*Plasma) IObject

type Plasma struct {
	currentId          int64
	mutex              *sync.Mutex
	IterStack          *IterStack
	builtInSymbolTable *SymbolTable
	BytecodeStack      *CodeStack
	MemoryStack        *ObjectStack
	TryStack           *TryStack
	SymbolTableStack   *SymbolStack
	Crc32Hash          hash.Hash32
	seed               uint64
	stdinScanner       *bufio.Scanner
	stdin              io.Reader
	stdout             io.Writer
	stderr             io.Writer
}

func (p *Plasma) NextId() int64 {
	p.mutex.Lock()
	result := p.currentId
	p.currentId++
	p.mutex.Unlock()
	return result
}

func (p *Plasma) PushObject(object IObject) {
	p.MemoryStack.Push(object)
}
func (p *Plasma) PeekObject() IObject {
	return p.MemoryStack.Peek()
}

func (p *Plasma) PopObject() IObject {
	return p.MemoryStack.Pop()
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
	p.builtInSymbolTable.Set(symbolName, loader(p))
}

func (p *Plasma) LoadBuiltInSymbols(symbolMap map[string]ObjectLoader) {
	for symbol, loader := range symbolMap {
		p.builtInSymbolTable.Set(symbol, loader(p))
	}
}

/*
	InitializeBytecode
	Loads the bytecode and clears the stack
*/
func (p *Plasma) InitializeBytecode(bytecode *Bytecode) {
	p.BytecodeStack.Clear()
	p.MemoryStack.Clear()
	p.SymbolTableStack.Clear()
	p.IterStack.Clear()
	p.TryStack.Clear()
	p.PushBytecode(bytecode)
	symbols := NewSymbolTable(p.builtInSymbolTable)
	symbols.Set("__built_in__",
		&Object{
			id:         p.NextId(),
			typeName:   ObjectName,
			class:      nil,
			subClasses: nil,
			symbols:    p.builtInSymbolTable,
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
	p.SymbolTableStack.Push(symbols)
}

func (p *Plasma) PushSymbolTable(table *SymbolTable) {
	p.SymbolTableStack.Push(table)
}

func (p *Plasma) PopSymbolTable() *SymbolTable {
	return p.SymbolTableStack.Pop()
}

func (p *Plasma) PeekSymbolTable() *SymbolTable {
	return p.SymbolTableStack.Peek()
}

func (p *Plasma) PushBytecode(code *Bytecode) {
	p.BytecodeStack.Push(code)
}

func (p *Plasma) PopBytecode() *Bytecode {
	return p.BytecodeStack.Pop()
}

func (p *Plasma) PeekBytecode() *Bytecode {
	return p.BytecodeStack.Peek()
}

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
	return p.builtInSymbolTable
}

func NewPlasmaVM(stdin io.Reader, stdout io.Writer, stderr io.Writer) *Plasma {
	number, randError := rand.Int(rand.Reader, big.NewInt(polySize))
	if randError != nil {
		panic(randError)
	}
	vm := &Plasma{
		currentId:        1,
		mutex:            &sync.Mutex{},
		IterStack:        NewIterStack(),
		BytecodeStack:    NewCodeStack(),
		MemoryStack:      NewObjectStack(),
		TryStack:         NewTryStack(),
		SymbolTableStack: NewSymbolStack(),
		Crc32Hash:        crc32.New(crc32.MakeTable(polySize)),
		seed:             number.Uint64(),
		stdinScanner:     bufio.NewScanner(stdin),
		stdin:            stdin,
		stdout:           stdout,
		stderr:           stderr,
	}
	vm.setBuiltInSymbols()
	return vm
}
