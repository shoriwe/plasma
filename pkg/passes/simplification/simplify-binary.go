package simplification

import (
	"github.com/shoriwe/gplasma/pkg/ast"
	"github.com/shoriwe/gplasma/pkg/common"
	"github.com/shoriwe/gplasma/pkg/lexer"
	"reflect"
)

func evalBinaryExpression(left, right *ast.BasicLiteralExpression, operator *lexer.Token) ast.IExpression {
	if reflect.TypeOf(left) != reflect.TypeOf(right) {
		return &ast.BinaryExpression{
			IExpression:   nil,
			LeftHandSide:  left,
			Operator:      operator,
			RightHandSide: right,
		}
	}

	switch operator.DirectValue {
	case lexer.Equals:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.SlicesEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.SlicesEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsString(left) && literalIsString(right):
			return booleanValue(
				common.SlicesEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsBytes(left) && literalIsBytes(right):
			return booleanValue(
				common.SlicesEqual(left.Token.Contents, right.Token.Contents),
			)
		}
	case lexer.NotEqual:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.SlicesNotEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.SlicesNotEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsString(left) && literalIsString(right):
			return booleanValue(
				common.SlicesNotEqual(left.Token.Contents, right.Token.Contents),
			)
		case literalIsBytes(left) && literalIsBytes(right):
			return booleanValue(
				common.SlicesNotEqual(left.Token.Contents, right.Token.Contents),
			)
		}
	case lexer.GreaterThan:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.IntegerGreaterToken(left.Token, right.Token),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.FloatGreaterToken(left.Token, right.Token),
			)
		}
	case lexer.GreaterOrEqualThan:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.IntegerGreatOrEqualToken(left.Token, right.Token),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.FloatGreatOrEqualToken(left.Token, right.Token),
			)
		}
	case lexer.LessThan:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.IntegerLessToken(left.Token, right.Token),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.FloatLessToken(left.Token, right.Token),
			)
		}
	case lexer.LessOrEqualThan:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			return booleanValue(
				common.IntegerLessOrEqualToken(left.Token, right.Token),
			)
		case literalIsFloat(left) && literalIsFloat(right):
			return booleanValue(
				common.FloatLessOrEqualToken(left.Token, right.Token),
			)
		}
	case lexer.BitwiseOr:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerBitwiseOrToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.BitwiseXor:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerBitwiseXorToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.BitWiseAnd:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerBitwiseAndToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.BitwiseLeft:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerBitwiseLeftToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.BitwiseRight:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerBitwiseRightToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.Add:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerAddToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatAddToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsString(left) && literalIsString(right):
			newValue := common.StringAddToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsBytes(left) && literalIsBytes(right):
			newValue := common.BytesStringAddToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.Sub:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerSubToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatSubToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.Star:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerMulToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatMulToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsString(left) && literalIsInteger(right):
			newValue := common.StringMulToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsInteger(left) && literalIsString(right):
			newValue := common.StringMulToken(right.Token, left.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsBytes(left) && literalIsInteger(right):
			newValue := common.BytesStringMulToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsInteger(left) && literalIsBytes(right):
			newValue := common.BytesStringMulToken(right.Token, left.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.Div:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerDivToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatDivToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.FloorDiv:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerDivToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatFloorDivToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.Modulus:
		switch {
		case literalIsInteger(left) && literalIsInteger(right):
			newValue := common.IntegerModToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	case lexer.PowerOf:
		switch {
		case literalIsFloat(left) && literalIsFloat(right):
			newValue := common.FloatPowToken(left.Token, right.Token)
			return &ast.BasicLiteralExpression{
				Token:       newValue,
				Kind:        newValue.Kind,
				DirectValue: newValue.DirectValue,
			}
		}
	}
	return &ast.BinaryExpression{
		LeftHandSide:  left,
		Operator:      operator,
		RightHandSide: right,
	}
}

func simplifyBinary(binaryExpression *ast.BinaryExpression) ast.IExpression {
	newLeft := SimplifyExpression(binaryExpression.LeftHandSide)
	newRight := SimplifyExpression(binaryExpression.RightHandSide)

	switch newLeft.(type) {
	case ast.BasicLiteralExpression:
		break
	default:
		return &ast.BinaryExpression{
			LeftHandSide:  newLeft,
			Operator:      binaryExpression.Operator,
			RightHandSide: newRight,
		}
	}
	switch newRight.(type) {
	case ast.BasicLiteralExpression:
		break
	default:
		return &ast.BinaryExpression{
			LeftHandSide:  newLeft,
			Operator:      binaryExpression.Operator,
			RightHandSide: newRight,
		}
	}
	return evalBinaryExpression(
		newLeft.(*ast.BasicLiteralExpression),
		newRight.(*ast.BasicLiteralExpression),
		binaryExpression.Operator,
	)
}
