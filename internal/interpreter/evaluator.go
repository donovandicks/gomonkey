package interpreter

import (
	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/object"
)

func evalProgram(stmts []ast.Statement) object.Object {
	var res object.Object

	for _, stmt := range stmts {
		res = Eval(stmt)

		switch res := res.(type) {
		case *object.ReturnVal:
			return res.Value
		case *object.Err:
			return res
		}
	}

	return res
}

func evalBlockStatement(block *ast.BlockStatement) object.Object {
	var res object.Object

	for _, stmt := range block.Statements {
		res = Eval(stmt)

		if res != nil {
			rt := res.Type()
			if rt == object.OBJ_ERR || rt == object.OBJ_RETURN {
				return res
			}
		}
	}

	return res
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
		return object.NewErr("invalid operator '-' for type %s", right.Type())
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
		return object.NewErr("unknown operator '%s' for type %s", operator, right.Type())
	}
}

func evalIntegerInfixExpr(operator string, left, right object.Object) object.Object {
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
		return object.NewErr("unknown integer operator '%s' on integers %d, %d", operator, l, r)
	}
}

func evalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJ_INTEGER && right.Type() == object.OBJ_INTEGER:
		return evalIntegerInfixExpr(operator, left, right)
	case operator == "==":
		return object.BoolFromNative(left == right)
	case operator == "!=":
		return object.BoolFromNative(left != right)
	case left.Type() != right.Type():
		return object.NewErr("type error: cannot perform '%s' on %s, %s", operator, left.Type(), right.Type())
	default:
		return object.NewErr("unknown operator '%s' for types %s, %s", operator, left.Type(), right.Type())
	}
}

func evalIfExpression(expr *ast.IfExpression) object.Object {
	cond := Eval(expr.Condition)
	if object.IsErr(cond) {
		return cond
	}

	if object.IsTruthy(cond) {
		return Eval(expr.Consequence)
	}

	if expr.Alternative != nil {
		return Eval(expr.Alternative)
	}

	return object.NullObject
}

func Eval(node ast.Node) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements)
	case *ast.ExpressionStatement:
		return Eval(node.Expression)
	case *ast.IntegerLiteral:
		return object.NewIntegerObject(node.Value)
	case *ast.Boolean:
		return object.BoolFromNative(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right) // evaluate the operand
		if object.IsErr(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.InfixExpression:
		left := Eval(node.Left)
		if object.IsErr(left) {
			return left
		}

		right := Eval(node.Right)
		if object.IsErr(right) {
			return right
		}

		return evalInfixExpr(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node)
	case *ast.IfExpression:
		return evalIfExpression(node)
	case *ast.ReturnStatement:
		val := Eval(node.Value)
		if object.IsErr(val) {
			return val
		}
		return object.NewReturnVal(val)
	}

	return nil
}
