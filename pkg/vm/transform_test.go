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

func TestPlasma_ToValueString(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.Symbols(), "Plasma")
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
}

func TestPlasma_ToValueBytes(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.Symbols(), []byte("Plasma"))
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
}

func TestPlasma_ToValueBool(t *testing.T) {
	p := NewVM(nil, nil, nil)
	s, err := p.ToValue(p.Symbols(), true)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, true, (<-rCh).Bool())
}

func TestPlasma_ToValueUint(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{uint(1), uintptr(1), uint8(1), uint16(1), uint32(1), uint64(1)} {
		s, err := p.ToValue(p.Symbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, int64(1), (<-rCh).Int())
	}
}

func TestPlasma_ToValueInt(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{int(1), int8(1), int16(1), int32(1), int64(1)} {
		s, err := p.ToValue(p.Symbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, int64(1), (<-rCh).Int())
	}
}

func TestPlasma_ToValueFloat(t *testing.T) {
	p := NewVM(nil, nil, nil)
	for _, element := range []any{float32(1), float64(1)} {
		s, err := p.ToValue(p.Symbols(), element)
		assert.Nil(t, err)
		p.Load("s", func(plasma *Plasma) *Value { return s })
		rCh, errCh, _ := p.ExecuteString("s")
		assert.Nil(t, <-errCh)
		assert.Equal(t, float64(1), (<-rCh).Float())
	}
}

func TestPlasma_ToValueSlice(t *testing.T) {
	p := NewVM(nil, nil, nil)
	slice := []any{1, 2, 3, "Plasma", "secret", 1.0}
	s, err := p.ToValue(p.Symbols(), slice)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s[0]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(1), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[1]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(2), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[2]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(3), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[3]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[4]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "secret", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[5]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1.0, (<-rCh).Float())
}

func TestPlasma_ToValueArray(t *testing.T) {
	p := NewVM(nil, nil, nil)
	slice := [6]any{1, 2, 3, "Plasma", "secret", 1.0}
	s, err := p.ToValue(p.Symbols(), slice)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s[0]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(1), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[1]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(2), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[2]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(3), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s[3]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "Plasma", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[4]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "secret", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s[5]")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 1.0, (<-rCh).Float())
}

func TestPlasma_ToValueMap(t *testing.T) {
	p := NewVM(nil, nil, nil)
	m := map[string]any{
		"john":  1,
		"conor": "Plasma",
	}
	s, err := p.ToValue(p.Symbols(), m)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s['john']")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(1), (<-rCh).Int())
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
	s, err := p.ToValue(p.Symbols(), obj)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s.Name")
	assert.Nil(t, <-errCh)
	assert.Equal(t, "sulcud", (<-rCh).String())
	rCh, errCh, _ = p.ExecuteString("s.Age")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(20), (<-rCh).Int())
	rCh, errCh, _ = p.ExecuteString("s.Health")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(100), (<-rCh).Int())
}

func TestPlasma_ToValueFunc(t *testing.T) {
	p := NewVM(nil, nil, nil)
	// No argument function
	f1 := func() int {
		return 1
	}
	s, err := p.ToValue(p.Symbols(), f1)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s()")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(f1()), (<-rCh).Int())
	// Arguments function
	f2 := func(a, b, c, d, e int) int {
		return a + b*c/d - e
	}
	s, err = p.ToValue(p.Symbols(), f2)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ = p.ExecuteString("s(1, 2, 3, 4, 5)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(f2(1, 2, 3, 4, 5)), (<-rCh).Int())
	// Function struct argument
	f3 := func(times int, ctx struct {
		Name string
	}) string {
		return strings.Repeat(ctx.Name, times)
	}
	s, err = p.ToValue(p.Symbols(), f3)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ = p.ExecuteString("ctx = Value()\nctx.Name = 'Plasma '\ns(3, ctx)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f3(3, struct{ Name string }{Name: "Plasma "}), (<-rCh).String())
	// Function struct pointer argument
	f4 := func(times int, ctx *struct {
		Name string
	}) string {
		return strings.Repeat(ctx.Name, times)
	}
	s, err = p.ToValue(p.Symbols(), f4)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ = p.ExecuteString("ctx = Value()\nctx.Name = 'Plasma '\ns(3, ctx)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, f4(3, &struct{ Name string }{Name: "Plasma "}), (<-rCh).String())
	// Function returns multiple values
	f5 := func(a int, b int) (int, int) {
		return a * 30, b * 3
	}
	s, err = p.ToValue(p.Symbols(), f5)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ = p.ExecuteString("s(3, 3)")
	assert.Nil(t, <-errCh)
	f5a1, f5a2 := f5(3, 3)
	result := <-rCh
	f5ca1 := result.GetValues()[0].Int()
	f5ca2 := result.GetValues()[1].Int()
	assert.Equal(t, int64(f5a1), f5ca1)
	assert.Equal(t, int64(f5a2), f5ca2)
}

func TestPlasma_ToValuePointer(t *testing.T) {
	p := NewVM(nil, nil, nil)
	a := 100
	b := &a
	c := &b
	s, err := p.ToValue(p.Symbols(), &c)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	rCh, errCh, _ := p.ExecuteString("s")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(100), (<-rCh).Int())
}

func TestPlasma_ToValueChan(t *testing.T) {
	p := NewVM(nil, nil, nil)
	a := make(chan int, 2)
	defer close(a)
	a <- 10
	s, err := p.ToValue(p.Symbols(), a)
	assert.Nil(t, err)
	p.Load("s", func(plasma *Plasma) *Value { return s })
	// Recv
	rCh, errCh, _ := p.ExecuteString("s.recv()")
	assert.Nil(t, <-errCh)
	assert.Equal(t, int64(10), (<-rCh).Int())
	// Send
	rCh, errCh, _ = p.ExecuteString("s.send(100)")
	assert.Nil(t, <-errCh)
	assert.Equal(t, 100, <-a)

}
