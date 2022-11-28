# Embedding

Embedding the language is simple:

```go
package main

import (
	"github.com/shoriwe/gplasma/pkg/vm"
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