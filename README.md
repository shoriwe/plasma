# gruby

## Description

**`gruby`** is a programming language made for the love and respect of the amazing **`ruby`** syntax and semantics with
a new features like interfaces to work with.

## Try it

You can have a working interpreter by compiling `cmd/plasma`.

1. Download the repository as a module
   
```shell
go get github.com/shoriwe/gruby
```

2. Then compile it

```shell
go build github.com/shoriwe/gruby/cmd/plasma -o plasma
```

```
...\Desktop> .\plasma.exe -h
...\Desktop\plasma.exe [FLAG [FLAG [FLAG]]] [PROGRAM [PROGRAM [PROGRAM]]]

[+] Notes
        - No PROGRAM arguments will spawn a REPL

[+] Flags
        -h, --help              Show this help message

[+] Environment Variables
        NoColor -> TRUE or FALSE                Disable color printing for this CLI
```

## Features

### Embedding

**`gruby`** was designed to be embedded in other go applications, you should do it like:

```go
package main

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/compiler/plasma"
	"github.com/shoriwe/gruby/pkg/reader"
	"github.com/shoriwe/gruby/pkg/vm"
	"os"
)

func main() {
	virtualMachine := vm.NewPlasmaVM(os.Stdin, os.Stdout, os.Stderr)
	compiler := plasma.NewCompiler(reader.NewStringReader("println('Hello world')"),
		map[uint8]uint8{
			plasma.PopRawExpressions: plasma.PopRawExpressions,
		},
	)
	code, compilationError := compiler.Compile()
	if compilationError != nil {
		panic(compilationError.String())
	}
	virtualMachine.InitializeByteCode(code)
	_, executionError := virtualMachine.Execute()
	if executionError != nil {
		panic(fmt.Sprintf("%s: %s", executionError.TypeName(), executionError.GetString()))
	}
}
```

In the future there will be a simpler way to embed it in your application, which shouldn't break the one provided
before.

### Language

This are the expressions and statements that **`gruby`** currently supports:

- [X] Expressions
    - [X] Literals
        - [X] Integers
            - [X] Decimal
            - [X] Hex
            - [X] Binary
            - [X] Octal
        - [X] Float
            - [X] Basic
            - [X] Scientific
        - [X] String
            - [X] Single Quote
            - [X] Double Quote
        - [X] Bytes
        - [X] Bool
            - [X] True
            - [X] False
        - [X] None
    - [X] Complex types
        - [X] Tuple
        - [X] Array
        - [X] Hash Table
    - [X] Unary Expressions
        - [X] Negate Bits
        - [X] Negate Bool
    - [X] Binary Expressions
        - [X] Add
        - [X] Sub
        - [X] Mul
        - [X] Div
        - [X] Mod
        - [X] Pow
        - [X] BitXor
        - [X] BitAnd
        - [X] BitOr
        - [X] BitLeft
        - [X] BitRight
        - [X] And
        - [X] Or
        - [X] Xor
        - [X] Equals
        - [X] NotEquals
        - [X] GreaterThan
        - [X] LessThan
        - [X] GreaterThanOrEqual
        - [X] LessThanOrEqual
    - [X] Lambda Expressions
    - [X] One Line If
    - [X] One Line Unless
    - [X] Identifiers
    - [X] Generators
    - [X] Call
        - [X] Function
        - [X] Type
    - [X] Index
    - [X] Parentheses Expressions
    - [X] Selector Expressions
- [ ] Statements
    - [X] Assign Statement
    - [ ] DeferStatement
    - [x] Do While
        - [X] Continue
        - [X] Break
        - [X] Redo
    - [X] While
        - [X] Continue
        - [X] Break
        - [X] Redo
    - [X] Until
        - [X] Continue
        - [X] Break
        - [X] Redo
    - [X] For Loop
        - [X] Continue
        - [X] Break
        - [X] Redo
    - [X] If - Else - Elif
    - [X] Unless - Else - Elif
    - [X] Switch
    - [X] Module
    - [X] Function Definition
    - [X] Interface
    - [X] Class
    - [X] Raise
    - [X] Try - Except
    - [X] Begin
    - [X] End
    - [X] Return
    - [ ] Super
    - [X] Pass
    - [ ] Yield

## Notable Differences

The major difference between **`ruby`** and **`gruby`** is that in the first the last expression in a function will be
returned without specifying the keyboard `return` but in **`gruby`** you should.

Another one will be that function calls, will always need parentheses to be executed, other way their will be evaluated
as objects.

Example:

This example shows a valid **`ruby`** code that returns from a function a string.

```ruby
def hello()
    "Hello World"
end

puts hello
```

But in **`gruby`** you should code it something like:

```ruby
def hello()
    return "Hello World" # Notice that here is used the keyboard "return"
end

println(hello())
```

## Important

- **`gruby`** is not a replacement for **`ruby`**, there are many missing features that are required to make **`ruby`**
  code run properly in **`gruby`**, of course there is the chance that some simpler scripts will run just fine
  in **`gruby`**.

- Also notice that this software has passed only by one iteration, so you should expect some bugs and performance
  issues, which will be of course corrected in the future.

# Useful references

This where useful references that made this project possible.

- [BNF grammar](https://ruby-doc.org/docs/ruby-doc-bundle/Manual/man-1.4/yacc.html)
- [Syntax Documentation](https://ruby-doc.org/docs/ruby-doc-bundle/Manual/man-1.4/syntax.html)