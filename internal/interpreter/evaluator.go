package interpreter

import (
	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/object"
)

var count = 0

func evalProgram(stmts []ast.Statement, env *object.Environment) object.Object {
	var res object.Object

	for _, stmt := range stmts {
		res = Eval(stmt, env)

		switch res := res.(type) {
		case *object.ReturnVal:
			return res.Value
		case *object.Err:
			return res
		}
	}

	return res
}

func evalBlockStatement(block *ast.BlockStatement, env *object.Environment) object.Object {
	var res object.Object

	for _, stmt := range block.Statements {
		res = Eval(stmt, env)

		if res != nil {
			rt := res.Type()
			if rt == object.OBJ_ERR || rt == object.OBJ_RETURN {
				return res
			}
		}
	}

	return res
}

func evalExpressions(exprs []ast.Expression, env *object.Environment) []object.Object {
	objs := make([]object.Object, 0, len(exprs))

	for _, expr := range exprs {
		val := Eval(expr, env)
		objs = append(objs, val)
		if object.IsErr(val) {
			return objs
		}
	}

	return objs
}

func evalIdentifier(node *ast.Identifier, env *object.Environment) object.Object {
	if ident, ok := env.Get(node.Value); ok {
		return ident
	}

	if builtin, ok := Builtins[node.Value]; ok {
		return builtin
	}

	return object.NewErr("undefined variable '%s'", node.Value)
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

func evalStringInfixExpr(operator string, left, right object.Object) object.Object {
	l := left.(*object.String).Value
	r := right.(*object.String).Value

	switch operator {
	case "+":
		return object.NewStringObject(l + r)
	default:
		return object.NewErr("unknown string operator '%s' on strings %s, %s", operator, l, r)
	}
}

func evalInfixExpr(operator string, left, right object.Object) object.Object {
	switch {
	case left.Type() == object.OBJ_INTEGER && right.Type() == object.OBJ_INTEGER:
		return evalIntegerInfixExpr(operator, left, right)
	case left.Type() == object.OBJ_STR && right.Type() == object.OBJ_STR:
		return evalStringInfixExpr(operator, left, right)
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

func evalIfExpression(expr *ast.IfExpression, env *object.Environment) object.Object {
	cond := Eval(expr.Condition, env)
	if object.IsErr(cond) {
		return cond
	}

	if object.IsTruthy(cond) {
		return Eval(expr.Consequence, env)
	}

	if expr.Alternative != nil {
		return Eval(expr.Alternative, env)
	}

	return object.NullObject
}

func evalWhileStatement(stmt *ast.WhileStatement, env *object.Environment) object.Object {
	for object.IsTruthy(Eval(stmt.Condition, env)) {
		evalBlockStatement(stmt.Block, env)
	}

	return object.NullObject
}

func unwrap(ret object.Object) object.Object {
	if r, ok := ret.(*object.ReturnVal); ok {
		return r.Value
	}

	return ret
}

func applyFunc(f object.Object, args []object.Object) object.Object {
	switch fn := f.(type) {
	case *object.Function:
		newEnv := object.NewEnvFromEnv(fn.Env)
		for idx, param := range fn.Parameters {
			newEnv.Set(param.Value, args[idx])
		}
		return unwrap(Eval(fn.Body, newEnv))
	case *object.Builtin:
		return fn.Fn(args...)
	default:
		return object.NewErr("undefined function '%s'", fn.Type())
	}

}

func Eval(node ast.Node, env *object.Environment) object.Object {
	switch node := node.(type) {
	case *ast.Program:
		return evalProgram(node.Statements, env)
	case *ast.ExpressionStatement:
		return Eval(node.Expression, env)
	case *ast.IntegerLiteral:
		return object.NewIntegerObject(node.Value)
	case *ast.StringLiteral:
		return object.NewStringObject(node.Value)
	case *ast.Boolean:
		return object.BoolFromNative(node.Value)
	case *ast.PrefixExpression:
		right := Eval(node.Right, env) // evaluate the operand
		if object.IsErr(right) {
			return right
		}
		return evalPrefixExpr(node.Operator, right)
	case *ast.AssignmentExpression:
		right := Eval(node.Right, env)
		if object.IsErr(right) {
			return right
		}

		return env.Update(node.Left.Value, right)
	case *ast.InfixExpression:
		left := Eval(node.Left, env)
		if object.IsErr(left) {
			return left
		}

		right := Eval(node.Right, env)
		if object.IsErr(right) {
			return right
		}

		return evalInfixExpr(node.Operator, left, right)
	case *ast.BlockStatement:
		return evalBlockStatement(node, env)
	case *ast.IfExpression:
		return evalIfExpression(node, env)
	case *ast.Identifier:
		return evalIdentifier(node, env)
	case *ast.WhileStatement:
		return evalWhileStatement(node, env)
	case *ast.LetStatement:
		val := Eval(node.Value, env)
		if object.IsErr(val) {
			return val
		}

		env.Set(node.Name.Value, val)
		return nil
	case *ast.FunctionLiteral:
		return object.NewFunctionObject(node.Parameters, node.Body, env)
	case *ast.CallExpression:
		f := Eval(node.Function, env)
		if object.IsErr(f) {
			return f
		}

		args := evalExpressions(node.Arguments, env)
		if len(args) > 0 && object.IsErr(args[0]) {
			return args[0]
		}

		return applyFunc(f, args)
	case *ast.ReturnStatement:
		val := Eval(node.Value, env)
		if object.IsErr(val) {
			return val
		}
		return object.NewReturnVal(val)
	}

	return nil
}
