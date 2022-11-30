# Embedding

## Basics concepts

- The plasma VM is thread safe by nature meaning you can have the same VM instance running different scripts at the same time. Accessing the same values at the same time.

- The VM has its own bytecode this was made this way to improve performance.

## Creating a new VM

To embed the language in your code use the function [NewVM](https://pkg.go.dev/github.com/shoriwe/plasma#NewVM)

```go
p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
```

If you are sure that your code doesn't need to print its output or read its input you can `nill` **stdin**, **stdout** and **stderr**.

```go
p := plasma.NewVM(nil, nil, nil)
```

## Executing your first script

To execute plasma functions you have two options, precompile them to the VM's bytecode or executing directly from a string. This two ways are completely different, for frequently running scripts is recommended the precompiled way since your code prepares the bytecode once, the string executing is only used for those scenarios where you **don't care** of the precompile time lost (which is small by the way).

### Executing a script from a string

This is the simples ways to execute plasma code, you only need to call [ExecuteString](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.ExecuteString) For example:

```go
p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
rCh, errCh, _ := p.ExecuteString("1 + 2")
err := <-errCh
if err != nil {
	panic(err)
}
fmt.Println(vm.Int[int](<-rCh))
```

### Precompiling a script

You will need to compile your script using [Compile](https://pkg.go.dev/github.com/shoriwe/plasma#Compile) then you can execute it with [Execute](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.Execute) For example:

```go
p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
bytecode, compileErr := plasma.Compile("1 + 2")
if compileErr != nil {
	panic(compileErr)
}
p.Execute(bytecode)
rCh, errCh, _ := p.Execute(bytecode)
err := <-errCh
if err != nil {
	panic(err)
}
fmt.Println(vm.Int[int](<-rCh))
```

### Why results of execution functions are channels?

As you have notice execution functions return channels, this was made to make use of the nature of thread safe execution to allow option to stop running scripts. You can stop a running script by sending an empty struct to the **stop channel** (Last return value of execution functions)

```go
p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
bytecode, compileErr := plasma.Compile("1 + 2")
if compileErr != nil {
	panic(compileErr)
}
p.Execute(bytecode)
rCh, errCh, stopCh := p.Execute(bytecode)
stopCh <- struct{}{} // Stop the running script
err := <-errCh
if err != nil {
	panic(err)
}
fmt.Println(vm.Int[int](<-rCh))
```

## Calling `Go` from `plasma`

There are two ways to call **Go** from **plasma**, by creating the functions manually or by letting the language convert everything for us.

### Automatic

Use the [LoadGo](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.LoadGo) function to transform any **Go** value to **plasma** values, this has clear limitations that will be mentioned in [Passing values from Go to `plasma`](#passing-values-from-go-to-plasma). In this example will pass a Go function directly to plasma:

```go
func Add(a, b int) int {
    return a + b
}

p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
loadErr := p.LoadGo("Add", Add)
if loadErr != nil {
	panic(loadErr)
}
_, errCh, _ := p.ExecuteString("println( Add(1, 2) )")
err := <-errCh
if err != nil {
	panic(err)
}
```

### Manual

Manual interfacing with Go is sometimes required for those scenarios the [LoadGo](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.LoadGo) function is unable to properly convert your **Go** code. In this kind of situation you will need to make use of the [Load](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.Load) function and [Loader](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Loader) interface. We will now recreate the example of before but using this method.

```go
func Add(a, b int) int {
	return a + b
}

// Function satisfying the Loader interface
func AddLoader(p *vm.Plasma) *vm.Value {
	addFunc := p.NewBuiltInFunction(
		p.RootSymbols(),
		func(argument ...*vm.Value) (*vm.Value, error) {
			a := vm.Int[int](argument[0])
			b := vm.Int[int](argument[1])
			return p.NewInt(int64(Add(a, b))), nil
		})
	return addFunc
}

p := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
p.Load("Add", AddLoader) // Loading the object
_, errCh, _ := p.ExecuteString("println( Add(1, 2) )")
err := <-errCh
if err != nil {
	panic(err)
}
```

## Passing values from `Go` to `plasma`

To speed up your interfacing with `plasma` you can make use of [LoadGo](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.LoadGo) to pass arbitrary Go values to the virtual machine. This has some limitations since it is still a feature in development but stable enough to resolve some scenarios. The current conversion table goes as follow.

| Go type                                                      | Plasma result          | Notes                                                        |
| ------------------------------------------------------------ | ---------------------- | ------------------------------------------------------------ |
| `int`, `int8`, `int16`, `int32`, `int64`, `uint`, `uintptr`, `uint8`, `uint16`, `uint32`, `uint64` | `Integer`              |                                                              |
| `float32`, `float64`                                         | `Float`                |                                                              |
| `string`                                                     | `String`               |                                                              |
| `bool`                                                       | `Bool`                 |                                                              |
| `complex64`, `comple64`                                      | Not supported yet      | Currently plasma doesn't support Go complex type             |
| Slices and Arrays                                            | `Array` or byte string | If the slice or array is of type `[]byte` or `[size]byte` it will be converted to a byte string |
| `map`                                                        | `Hash`                 | If the key or value type is still not supported it will fail to convert the entire map |
| Structs                                                      | `Value`                | Structs will be converted to `Value` objects, with all possible public fields of it, including struct methods and fields. Notice that it is recommended to use pointer structs (`&Struct`) instead of direct values |
| Functions                                                    | `BuiltInFunction`      |                                                              |
| Pointers, `unsafe.Pointer`                                   |                        | Pointers first resolve to the targeted pointed value the transform it to plasma objects |
| Channels                                                     | `Value`                | Channels are converted to `Value` objects with two specials methods. **`recv`** which internally does the **`<-channel`** operation and send **`send(VALUE_ARGUMENT)`** which internally does the **`channel <- VALUE_ARGUMENT`** |
| Interface                                                    | Not supported yet      | Interfaces are intended to be supported but not yet          |

If you want to convert `plasma` values to go values you can make use of [FromValue](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.FromValue). This function is able to convert any `plasma` value except values of these types: `BuiltInFunction`, `Function`, `BuiltInClass`, `Class`

## Working example

### main.go

```go
package main

import (
	"os"

	"github.com/shoriwe/plasma/pkg/vm"
)

const myScript = `
args = get_args()
if args.__len__() > 1
    println(args.__string__())
else
    println("No")
end
`

func main() {
	p := vm.NewVM(os.Stdin, os.Stdout, os.Stderr)
	p.Load("get_args", func(plasma *vm.Plasma) *vm.Value {
		return plasma.NewBuiltInFunction(p.RootSymbols(),
			func(argument ...*vm.Value) (*vm.Value, error) {
				tupleValues := make([]*vm.Value, 0, len(os.Args))
				for _, cmdArg := range os.Args {
					tupleValues = append(tupleValues, plasma.NewString([]byte(cmdArg)))
				}
				return plasma.NewTuple(tupleValues), nil
			})
	})
	_, errorChannel, _ := p.ExecuteString(myScript)
	err := <-errorChannel
	if err != nil {
		panic(err)
	}
}

```

Run this program with

```shell
go run main.go 1 2 3 4
```

