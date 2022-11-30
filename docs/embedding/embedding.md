# Embedding

## Basics

To embed the language in your code use the function [NewVM](https://pkg.go.dev/github.com/shoriwe/plasma#NewVM)

```go
vm := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
```

If you are sure that your code doesn't need to print its output or read its input you can `nill` **stdin**, **stdout** and **stderr**.

```go
vm := plasma.NewVM(nil, nil, nil)
```

## Calling `Go` from `plasma`

There are two ways to call **Go** from **plasma**, by creating the functions manually or by letting the language convert everything for us.

### Automatic

Use the [LoadGo](https://pkg.go.dev/github.com/shoriwe/plasma/pkg/vm#Plasma.LoadGo) function to transform any Go value to the **plasma** values, this has clear limitations that will be mentioned in [Passing values from Go to `plasma`](#passing-values-from-go-to-plasma)

In this example will pass a Go function directly to plasma

```go
func Add(a, b int) int {
    return a + b
}

vm := plasma.NewVM(os.Stdin, os.Stdout, os.Stderr)
plasmaAdd, transformErr := vm.ToValue(Add)
if transformErr != nil {
    panic(transformErr)
}
vm.LoadGo("")
```



### Manual

## Passing values from Go to `plasma`

## Why

```go
package main

import (
	"github.com/shoriwe/plasma/pkg/vm"
	"os"
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
	plasma := vm.NewVM(os.Stdin, os.Stdout, os.Stderr)
	plasma.Load("get_args", func(plasma *vm.Plasma) *vm.Value {
		return plasma.NewBuiltInFunction(plasma.Symbols(),
			func(argument ...*vm.Value) (*vm.Value, error) {
				tupleValues := make([]*vm.Value, 0, len(os.Args))
				for _, cmdArg := range os.Args {
					tupleValues = append(tupleValues, plasma.NewString([]byte(cmdArg)))
				}
				return plasma.NewTuple(tupleValues), nil
			})
	})
	_, errorChannel, _ := plasma.ExecuteString(myScript)
	err := <-errorChannel
	if err != nil {
		panic(err)
	}
}
```