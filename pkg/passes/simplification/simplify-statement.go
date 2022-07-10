package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/ast2"
)

func (simp *simplify) simplifyStatement(stmt ast.Statement) ast2.Statement {
	switch s := stmt.(type) {
	case *ast.AssignStatement:
		return simp.simplifyAssign(s)
	case *ast.DoWhileStatement:
		return simp.simplifyDoWhile(s)
	case *ast.WhileLoopStatement:
		return simp.simplifyWhile(s)
	case *ast.UntilLoopStatement:
		return simp.simplifyUntil(s)
	case *ast.ForLoopStatement:
		return simp.simplifyFor(s)
	case *ast.IfStatement:
		return simp.simplifyIf(s)
	case *ast.UnlessStatement:
		return simp.simplifyUnless(s)
	case *ast.SwitchStatement:
		return simp.simplifySwitch(s)
	case *ast.ModuleStatement:
		return simp.simplifyModule(s)
	case *ast.FunctionDefinitionStatement:
		return simp.simplifyFunction(s)
	case *ast.GeneratorDefinitionStatement:
		return simp.simplifyGeneratorDef(s)
	case *ast.InterfaceStatement:
		return simp.simplifyInterface(s)
	case *ast.ClassStatement:
		return simp.simplifyClass(s)
	case *ast.ReturnStatement:
		return simp.simplifyReturn(s)
	case *ast.YieldStatement:
		return simp.simplifyYield(s)
	case *ast.ContinueStatement:
		return simp.simplifyContinue(s)
	case *ast.BreakStatement:
		return simp.simplifyBreak(s)
	case *ast.PassStatement:
		return simp.simplifyPass(s)
	case *ast.BlockStatement:
		return simp.simplifyBlock(s)
	case *ast.RequireStatement:
		return simp.simplifyRequire(s)
	case *ast.DeleteStatement:
		return simp.simplifyDelete(s)
	case *ast.DeferStatement:
		return simp.simplifyDefer(s)
	default:
		panic("unknown statement type")
	}
}
