package transformations_1

import (
	"fmt"
	"github.com/shoriwe/gplasma/pkg/ast2"
	"github.com/shoriwe/gplasma/pkg/ast3"
	"reflect"
)

func getRoot(expr ast3.Expression) ast3.Expression {
	current := expr
	for {
		switch c := current.(type) {
		case *ast3.Selector:
			current = c.X
		case *ast3.Index:
			current = c.Source
		default:
			return c
		}
	}
}

type generatorTransform struct {
	selfSymbols map[string]struct{} // TODO: Init me
}

func (gt *generatorTransform) resolve(node ast3.Node, symbols map[string]struct{}) ast3.Node {

}

// Enumerate self symbols
func (gt *generatorTransform) enumerate(a ast3.Assignable) {
	switch left := a.(type) {
	case *ast3.Identifier:
		gt.selfSymbols[left.Symbol] = struct{}{}
	case *ast3.Selector:
		root := getRoot(left.X)
		ident, isIdent := root.(*ast3.Identifier)
		if !isIdent {
			return
		}
		if _, found := gt.selfSymbols[ident.Symbol]; !found {
			return
		}
		gt.selfSymbols[ident.Symbol] = struct{}{}
	case *ast3.Index:
		root := getRoot(left.Source)
		ident, isIdent := root.(*ast3.Identifier)
		if !isIdent {
			return
		}
		if _, found := gt.selfSymbols[ident.Symbol]; !found {
			return
		}
		gt.selfSymbols[ident.Symbol] = struct{}{}
	default:
		panic(fmt.Sprintf("unknown assignable type %s", reflect.TypeOf(a).String()))
	}
}

/*
	- Transform assignments to self.IDENTIFIER
	- Update Identifier access to self.IDENTIFIER
	- Add update jump condition variable before yield
	- Add label after yield
	- Add update jump condition variable before return
*/
func (gt *generatorTransform) process(node ast3.Node) []ast3.Node {
	switch n := node.(type) {

	}
}

/*
	- Prepend jump table
	- Append On finish label
	- Append Return none
*/
func (gt *generatorTransform) setup(body []ast3.Node) []ast3.Node {
	
}

func (gt *generatorTransform) transform(rawBody []ast3.Node) []ast3.Node {
	body := make([]ast3.Node, 0, len(rawBody))
	for _, node := range rawBody {
		body = append(body, gt.process(node)...)
	}
	return gt.setup(body)
}

func (transform *transformPass) GeneratorDef(generator *ast2.GeneratorDefinition) []ast3.Node {
	rawBody := make([]ast3.Node, 0, len(generator.Body))
	for _, node := range generator.Body {
		rawBody = append(rawBody, transform.Node(node)...)
	}
	gt := generatorTransform{selfSymbols: map[string]struct{}{}}
	return gt.transform(rawBody)
}
