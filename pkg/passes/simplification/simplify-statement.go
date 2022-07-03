package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"reflect"
)

func SimplifyStatement(statement ast.IStatement) ast.IStatement {
	switch stmt := statement.(type) {
	case *ast.AssignStatement:
		return simplifyAssign(stmt)
	case *ast.DoWhileStatement:
		return simplifyDoWhile(stmt)
	case *ast.WhileLoopStatement:
		return simplifyWhile(stmt)
	case *ast.UntilLoopStatement:
		return simplifyUntil(stmt)
	case *ast.ForLoopStatement:
		return simplifyForLoop(stmt)
	case *ast.SwitchStatement:
		return simplifySwitch(stmt)
	case *ast.ModuleStatement:
		return simplifyModule(stmt)
	case *ast.FunctionDefinitionStatement:
		return simplifyFunctionDef(stmt)
	case *ast.InterfaceStatement:
		return simplifyInterface(stmt)
	case *ast.ClassStatement:
		return simplifyClass(stmt)
	case *ast.TryStatement:
		return simplifyTry(stmt)
	case *ast.BeginStatement:
		return simplifyBegin(stmt)
	case *ast.EndStatement:
		return simplifyEnd(stmt)
	case *ast.ReturnStatement:
		return simplifyEnd(stmt)
	case *ast.YieldStatement:
		return simplifyYield(stmt)
	case *ast.RaiseStatement:
		return simplifyRaise(statement.X)
	case *ast.IfStatement:
		return simplifyIf(stmt)
	case *ast.UnlessStatement:
		return simplifyUnless(stmt)
	default:
		panic(reflect.TypeOf(stmt))
	}
}
