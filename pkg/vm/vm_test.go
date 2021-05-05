package vm

import (
	"fmt"
	"github.com/shoriwe/gruby/pkg/vm/object"
	_ "github.com/shoriwe/gruby/pkg/vm/object"
	_ "github.com/shoriwe/gruby/pkg/vm/utils"
	"testing"
)

func TestData(t *testing.T) {
	fmt.Println(object.NewString(nil, nil, "Hello").Id())
}
