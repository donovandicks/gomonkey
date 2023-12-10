package interpreter

import (
	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/object"
)

func evalStatements(stmts []ast.Statement) object.Object {
	var result object.Object

	for _, stmt := range stmts {
		result = Eval(stmt)
	}

	return result
}

func evalBangOpExpr(right object.Object) object.Object {
	switch right {
	case object.TrueBool:
		return object.FalseBool
	case object.FalseBool:
		return object.TrueBool
	case object.NullObject:
		return object.TrueBool
	default:
		return object.FalseBool
	}
}

func evalMinusOpExpr(right object.Object) object.Object {
	if right.Type() != object.OBJ_INTEGER {
		return object.NullObject
	}

	val := right.(*object.Integer).Value
	return object.NewIntegerObject(-val)
}

func evalPrefixExpr(operator string, right object.Object) object.Object {
	switch operator {
	case "!":
		return evalBangOpExpr(right)
	case "-":
		return evalMinusOpExpr(right)
	default:
		return object.NullObject
	}
}

func evalIntegerInfixExpr(operator string, left object.Object, right object.Object) object.Object {
	l := left.(*object.Integer).Value
	r := right.(*object.Integer).Value

	switch operator {
	case "+":
		return object.NewIntegerObject(l + r)
	case "-":
		return object.NewIntegerObject(l - r)
	case "*":
		return object.NewIntegerObject(l * r)
	case "/":
		return object.NewIntegerObject(l / r)
	case "<":
		return object.BoolFromNative(l < r)
	case ">":
		return object.BoolFromNative(l > r)
	case "==":
		return object.BoolFromNative(l == r)
	case "!=":
		return object.BoolFromNative(l != r)
	default:
		return object.NullObject
	}
}

func evalInfixExpr(operator string, left object.Object, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJ_INTEGER && right.Type() == object.OBJ_INTEGER:
		return evalIntegerInfixExpr(operator, left, right)
	case operator == "==":
		return object.BoolFromNative(left == right)
	case operator == "!=":
		return object.BoolFromNative(left != right)
	default:
		return object.NullObject
	}
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalStatements(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return object.NewIntegerObject(node.Value)
	case *ast.Boolean:
		return object.BoolFromNative(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right) // evaluate the operand
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		right := Eval(node.Right)
		return evalInfixExpr(node.Operator, left, right)
	}

	return nil
}
