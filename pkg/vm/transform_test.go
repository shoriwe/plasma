package vm

import (
	"github.com/stretchr/testify/assert"
	"strings"
	"testing"
)

func TestPlasma_FromValueValue(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("a = Value()\na.Name = \"Plasma\"\na")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, "Plasma", s.(map[string]any)["Name"])
}

func TestPlasma_FromValueString(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("'Plasma'")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, "Plasma", s.(string))
}

func TestPlasma_FromValueBytes(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("b\"Plasma\"")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, []byte("Plasma"), s)
}

func TestPlasma_FromValueBool(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("true")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, true, s)
}

func TestPlasma_FromValueNone(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("none")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, nil, s)
}

func TestPlasma_FromValueInt(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("10")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, int64(10), s)
}

func TestPlasma_FromValueFloat(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("10.0")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	assert.Equal(t, float64(10.0), s)
}

func TestPlasma_FromValueArray(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("[10.0, 20, '30']")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	ref := []any{10.0, int64(20), "30"}
	for index, v := range s.([]any) {
		assert.Equal(t, ref[index], v)
	}
}

func TestPlasma_FromValueMap(t *testing.T) {
	p := NewVM(nil, nil, nil)
	rCh, errCh, _ := p.ExecuteString("{'Hello': 1, 1: 10, 65.5: 0}")
	assert.Nil(t, <-errCh)
	s, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	ref := map[any]any{"Hello": int64(1), int64(1): int64(10), 65.5: int64(0)}
	ss := s.(map[any]any)
	for key, value := range ref {
		assert.Equal(t, value, ss[key])
	}
}

func TestPlasma_ToValueString(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.RootSymbols(), "Plasma")
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
}

func TestPlasma_ToValueBytes(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.RootSymbols(), []byte("Plasma"))
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
}

func TestPlasma_ToValueBool(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.RootSymbols(), true)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, true, (<-rCh).Bool())
}

func TestPlasma_ToValueUint(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{uint(1), uintptr(1), uint8(1), uint16(1), uint32(1), uint64(1)} {
		s, err := p.ToValue(p.RootSymbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, 1, Int[int](<-rCh))
	}
}

func TestPlasma_ToValueInt(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{int(1), int8(1), int16(1), int32(1), int64(1)} {
		s, err := p.ToValue(p.RootSymbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, 1, Int[int](<-rCh))
	}
}

func TestPlasma_ToValueFloat(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{float32(1), float64(1)} {
		s, err := p.ToValue(p.RootSymbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, float64(1), Float[float64](<-rCh))
	}
}

func TestPlasma_ToValueSlice(t *testing.T) {
	p := NewVM(nil, nil, nil)
	slice := []any{1, 2, 3, "Plasma", "secret", 1.0}
	s, err := p.ToValue(p.RootSymbols(), slice)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s[0]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[1]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 2, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[2]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 3, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[3]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[4]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "secret", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[5]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1.0, Float[float64](<-rCh))
}

func TestPlasma_ToValueArray(t *testing.T) {
	p := NewVM(nil, nil, nil)
	slice := [6]any{1, 2, 3, "Plasma", "secret", 1.0}
	s, err := p.ToValue(p.RootSymbols(), slice)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s[0]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[1]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 2, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[2]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 3, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s[3]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[4]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "secret", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[5]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1.0, Float[float64](<-rCh))
}

func TestPlasma_ToValueMap(t *testing.T) {
	p := NewVM(nil, nil, nil)
	m := map[string]any{
		"john":  1,
		"conor": "Plasma",
	}
	s, err := p.ToValue(p.RootSymbols(), m)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s['john']")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s['conor']")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
}

func TestPlasma_ToValueStruct(t *testing.T) {
	p := NewVM(nil, nil, nil)
	obj := struct {
		Name   string
		Age    int
		Health int
	}{
		"sulcud",
		20,
		100,
	}
	s, err := p.ToValue(p.RootSymbols(), obj)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s.Name")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "sulcud", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s.Age")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 20, Int[int](<-rCh))
	rCh, errCh, _ = p.ExecuteString("s.Health")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 100, Int[int](<-rCh))
}

func TestPlasma_ToValueFuncNoArgs(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// No argument function
	f := func() int {
		return 1
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s()")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(), Int[int](<-rCh))
}

func TestPlasma_ToValueFuncBasicType(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// Arguments function
	f := func(a, b, c, d, e int) int {
		return a + b*c/d - e
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s(1, 2, 3, 4, 5)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(1, 2, 3, 4, 5), Int[int](<-rCh))
}

func TestPlasma_ToValueFuncStructType(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// Function struct argument
	f := func(times int, ctx struct {
		Name string
	}) string {
		return strings.Repeat(ctx.Name, times)
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("ctx = Value()\nctx.Name = 'Plasma '\ns(3, ctx)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(3, struct{ Name string }{Name: "Plasma "}), (<-rCh).String())
}

func TestPlasma_ToValueStructPointer(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// Function struct pointer argument
	f := func(times int, ctx *struct {
		Name string
	}) string {
		return strings.Repeat(ctx.Name, times)
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("ctx = Value()\nctx.Name = 'Plasma '\ns(3, ctx)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(3, &struct{ Name string }{Name: "Plasma "}), (<-rCh).String())
}

func TestPlasma_ToValueFuncTypeAlias(t *testing.T) {
	type stringAlias string
	p := NewVM(nil, nil, nil)
	// Function type alias argument
	f := func(times int, s stringAlias) string {
		return strings.Repeat(string(s), times)
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s(3, 'Plasma ')")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(3, "Plasma "), (<-rCh).String())
}

func TestPlasma_ToValueFuncVariadicBasicType(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// Function type alias argument
	f := func(words ...string) string {
		return strings.Join(words, " ")
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s('Plasma', 'Plasma', 'Plasma')")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f("Plasma", "Plasma", "Plasma"), (<-rCh).String())
}

func TestPlasma_ToValueFuncVariadicTypeAlias(t *testing.T) {
	p := NewVM(nil, nil, nil)
	type stringAlias string
	// Function type alias argument
	f := func(words ...stringAlias) string {
		a := make([]string, 0, len(words))
		for _, b := range words {
			a = append(a, string(b))
		}
		return strings.Join(a, " ")
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s('Plasma', 'Plasma', 'Plasma')")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f("Plasma", "Plasma", "Plasma"), (<-rCh).String())
}

func TestPlasma_ToValueFuncVariadicStruct(t *testing.T) {
	p := NewVM(nil, nil, nil)
	type tt struct {
		Name string
	}
	// Function type alias argument
	f := func(words ...tt) string {
		a := make([]string, 0, len(words))
		for _, b := range words {
			a = append(a, b.Name)
		}
		return strings.Join(a, " ")
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("v = Value()\nv.Name = 'Plasma'\ns(v, v, v)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(tt{"Plasma"}, tt{"Plasma"}, tt{"Plasma"}), (<-rCh).String())
}

type testStruct struct {
	data string
}

func (ts *testStruct) Ptr() (string, string) {
	return "called from pointer", ts.data
}

func (ts testStruct) Value() (string, string) {
	return "called from value", ts.data
}

func TestPlasma_ToValueStructFunctions(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// When Pointer
	s, err := p.ToValue(p.RootSymbols(), &testStruct{"yes"})
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s.Ptr()")
	assert.Nil(t, <-errCh)
	gV, err := p.FromValue(<-rCh)
	assert.Nil(t, err)
	v := gV.([]any)
	assert.Equal(t, 2, len(v))
	assert.Equal(t, "called from pointer", v[0])
	assert.Equal(t, "yes", v[1])
	rCh, errCh, _ = p.ExecuteString("s.Value()")
	assert.Nil(t, <-errCh)
	gV, err = p.FromValue(<-rCh)
	assert.Nil(t, err)
	v = gV.([]any)
	assert.Equal(t, 2, len(v))
	assert.Equal(t, "called from value", v[0])
	assert.Equal(t, "yes", v[1])
	// When Value
	s, err = p.ToValue(p.RootSymbols(), testStruct{"yes"})
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	_, errCh, _ = p.ExecuteString("s.Ptr()")
	assert.NotNil(t, <-errCh)
	rCh, errCh, _ = p.ExecuteString("s.Value()")
	assert.Nil(t, <-errCh)
	gV, err = p.FromValue(<-rCh)
	assert.Nil(t, err)
	v = gV.([]any)
	assert.Equal(t, 2, len(v))
	assert.Equal(t, "called from value", v[0])
	assert.Equal(t, "yes", v[1])
}

func TestPlasma_ToValueFuncVariadicStructPtr(t *testing.T) {
	p := NewVM(nil, nil, nil)
	type tt struct {
		Name string
	}
	// Function type alias argument
	f := func(words ...*tt) string {
		a := make([]string, 0, len(words))
		for _, b := range words {
			a = append(a, b.Name)
		}
		return strings.Join(a, " ")
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("v = Value()\nv.Name = 'Plasma'\ns(v, v, v)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f(&tt{"Plasma"}, &tt{"Plasma"}, &tt{"Plasma"}), (<-rCh).String())
}

func TestPlasma_ToValueFuncMultipleReturnValues(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// Function returns multiple values
	f := func(a int, b int) (int, int) {
		return a * 30, b * 3
	}
	s, err := p.ToValue(p.RootSymbols(), f)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s(3, 3)")
	assert.Nil(t, <-errCh)
	a1, a2 := f(3, 3)
	result := <-rCh
	ca1 := Int[int](result.GetValues()[0])
	ca2 := Int[int](result.GetValues()[1])
	assert.Equal(t, a1, ca1)
	assert.Equal(t, a2, ca2)
}

type testPlasma_ToValueAliasMemberFunc_stringAlias string

func (s testPlasma_ToValueAliasMemberFunc_stringAlias) Say(word string) string {
	return string(s) + word
}

func TestPlasma_ToValueAliasMemberFunc(t *testing.T) {
	p := NewVM(nil, nil, nil)
	tt := testPlasma_ToValueAliasMemberFunc_stringAlias("Say ")
	s, err := p.ToValue(p.RootSymbols(), tt)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s.Say('Plasma')")
	assert.Nil(t, <-errCh)
	assert.Equal(t, tt.Say("Plasma"), (<-rCh).String())
}

type testPlasma_ToValueAliasMemberFuncVariadic string

func (s testPlasma_ToValueAliasMemberFuncVariadic) Say(word string, b ...string) string {
	for _, ss := range b {
		word += ss
	}
	return string(s) + word
}

func TestPlasma_ToValueAliasMemberFuncVariadic(t *testing.T) {
	p := NewVM(nil, nil, nil)
	tt := testPlasma_ToValueAliasMemberFuncVariadic("Say ")
	s, err := p.ToValue(p.RootSymbols(), tt)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	// Say
	rCh, errCh, _ := p.ExecuteString("s.Say('Plasma', 'hello', 'Plasma')")
	assert.Nil(t, <-errCh)
	assert.Equal(t, tt.Say("Plasma", "hello", "Plasma"), (<-rCh).String())
}

func TestPlasma_ToValuePointer(t *testing.T) {
	p := NewVM(nil, nil, nil)
	a := 100
	b := &a
	c := &b
	s, err := p.ToValue(p.RootSymbols(), &c)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 100, Int[int](<-rCh))
}

func TestPlasma_ToValueChan(t *testing.T) {
	p := NewVM(nil, nil, nil)
	a := make(chan int, 2)
	defer close(a)
	a <- 10
	s, err := p.ToValue(p.RootSymbols(), a)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	// Recv
	rCh, errCh, _ := p.ExecuteString("s.recv()")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 10, Int[int](<-rCh))
	// Send
	rCh, errCh, _ = p.ExecuteString("s.send(100)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 100, <-a)

}
