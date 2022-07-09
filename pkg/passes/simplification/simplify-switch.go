package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func simplifySwitch(switch_ *ast.SwitchStatement) *ast2.If {
	anonymousIdentifier := nextAnonIdentifier()
	root := &ast2.If{
		SwitchSetup: &ast2.Assignment{
			Left:  anonymousIdentifier,
			Right: simplifyExpression(switch_.Target),
		},
		Condition: nil,
		Body:      nil,
		Else:      nil,
	}
	currentIf := root
	for caseIndex, case_ := range switch_.CaseBlocks {
		for targetIndex, caseTarget := range case_.Cases {
			if targetIndex == 0 {
				currentIf.Condition = &ast2.Binary{
					Left:     anonymousIdentifier,
					Right:    simplifyExpression(caseTarget),
					Operator: ast2.Equals,
				}
				continue
			}
			currentIf.Condition = &ast2.Binary{
				Left: currentIf.Condition,
				Right: &ast2.Binary{
					Left:     anonymousIdentifier,
					Right:    simplifyExpression(caseTarget),
					Operator: ast2.Equals,
				},
				Operator: ast2.Or,
			}
		}
		currentIf.Body = make([]ast2.Node, 0, len(case_.Body))
		for _, node := range case_.Body {
			currentIf.Body = append(currentIf.Body, simplifyNode(node))
		}
		if caseIndex+1 < len(switch_.CaseBlocks) {
			newCurrentIf := &ast2.If{
				Condition: nil,
				Body:      nil,
				Else:      nil,
			}
			currentIf.Else = []ast2.Node{newCurrentIf}
			currentIf = newCurrentIf
		}
	}
	currentIf.Else = make([]ast2.Node, 0, len(switch_.Default))
	for _, node := range switch_.Default {
		currentIf.Else = append(currentIf.Else, simplifyNode(node))
	}
	return root
}
