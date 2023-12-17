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
		return object.NewErr(
			"type error: cannot perform '%s' on %s, %s",
			operator,
			left.Type(),
			right.Type(),
		)
	default:
		return object.NewErr(
			"unknown operator '%s' for types %s, %s",
			operator,
			left.Type(),
			right.Type(),
		)
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

func evalMapLiteral(node *ast.MapLiteral, env *object.Environment) object.Object {
	pairs := make(map[object.HashKey]object.KVPair, len(node.Entries))
	for key, val := range node.Entries {
		k := Eval(key, env)
		if object.IsErr(k) {
			return k
		}

		hashable, ok := k.(object.HashableObject)
		if !ok {
			return object.NewErr("cannot use unhashable type %s as hash key", k.Type())
		}

		v := Eval(val, env)
		if object.IsErr(v) {
			return v
		}

		pairs[hashable.Hash()] = object.KVPair{Key: k, Value: v}
	}

	return &object.Map{Entries: pairs}
}

func evalListIndexExpr(left, index object.Object) object.Object {
	l := left.(*object.List)
	idx := index.(*object.Integer).Value

	max := int64(len(l.Elems) - 1)
	if idx > max {
		return object.NewErr("index out of bounds: %d", idx)
	}

	if idx < 0 {
		idx = int64(len(l.Elems)) + idx
	}

	return l.Elems[idx]
}

func evalMapIndexExpr(left, index object.Object) object.Object {
	l := left.(*object.Map)
	key, _ := index.(object.HashableObject)

	kv, ok := l.Entries[key.Hash()]
	if !ok {
		return object.NewErr("no key found for %s (hash=%d)", key.Inspect(), key.Hash().Value)
	}

	return kv.Value
}

func evalIndexExpr(left, index object.Object) object.Object {
	switch left.Type() {
	case object.OBJ_LIST:
		if index.Type() != object.OBJ_INTEGER {
			return object.NewErr("cannot index list using non-integer type %s", left.Type())
		}
		return evalListIndexExpr(left, index)
	case object.OBJ_MAP:
		if !object.IsHashable(index) {
			return object.NewErr("cannot index map using non-hashable type %s", left.Type())
		}
		return evalMapIndexExpr(left, index)
	default:
		return object.NewErr("cannot index %s object", left.Type())
	}
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

func applyFunc(callable object.Object, args []object.Object) object.Object {
	switch c := callable.(type) {
	case *object.Class:
		return object.NewInstance(c)
	case *object.Function:
		newEnv := object.NewEnvFromEnv(c.Env)
		for idx, param := range c.Parameters {
			newEnv.Set(param.Value, args[idx])
		}
		return unwrap(Eval(c.Body, newEnv))
	case *object.Builtin:
		return c.Fn(args...)
	default:
		return object.NewErr("undefined callable '%s'", c.Type())
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
	case *ast.ListLiteral:
		elems := evalExpressions(node.Elems, env)
		if len(elems) >= 1 && object.IsErr(elems[0]) {
			return elems[0]
		}

		return object.NewListObject(elems)
	case *ast.MapLiteral:
		return evalMapLiteral(node, env)
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

		switch left := node.Left.(type) {
		case *ast.Identifier:
			return env.Update(left.Value, right)
		case *ast.GetExpression:
			evaled := Eval(left.Left, env)
			inst, ok := evaled.(*object.Instance)
			if !ok {
				return object.NewErr("cannot assign field to non-instance type %T", evaled)
			}

			field, ok := left.Right.(*ast.Identifier)
			if !ok {
				return object.NewErr("undefined property %s on %s", left.Right.String(), inst.Inspect())
			}

			inst.Set(field.String(), right)
			return right
		default:
			object.NewErr("cannot assign to %s (%T)", left, left)
		}
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
	case *ast.IndexExpression:
		left := Eval(node.Left, env)
		if object.IsErr(left) {
			return left
		}

		index := Eval(node.Index, env)
		if object.IsErr(index) {
			return index
		}

		return evalIndexExpr(left, index)
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
	case *ast.FunctionStatement: // named func stmt
		fn := object.NewFunctionObject(node.Name, node.Parameters, node.Body, env)
		env.Set(node.Name.Value, fn)
		return nil
	case *ast.FunctionLiteral: // anon func expr
		return object.NewFunctionObject(nil, node.Parameters, node.Body, env)
	case *ast.ClassStatement:
		cls := object.NewClassObject(node.Name, node.Methods, env)
		env.Set(node.Name.Value, cls)
		return nil
	case *ast.GetExpression:
		obj := Eval(node.Left, env)
		inst, ok := obj.(*object.Instance)
		if !ok {
			return object.NewErr("object %s has no properties", obj.Type())
		}

		val := inst.Get(node.Right.String())
		if val == nil {
			return object.NewErr("object %s has no property %s", inst.Inspect(), node.Right.String())
		}

		return val
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
