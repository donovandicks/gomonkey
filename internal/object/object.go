package object

import (
	"fmt"
	"strings"

	"github.com/donovandicks/gomonkey/internal/ast"
)

type ObjectType string
type BuiltinFn func(args ...Object) Object

const (
	OBJ_INTEGER ObjectType = "INTEGER"
	OBJ_BOOLEAN ObjectType = "BOOLEAN"
	OBJ_FUNC    ObjectType = "FUNCTION"
	OBJ_NULL    ObjectType = "NULL"
	OBJ_RETURN  ObjectType = "RETURN"
	OBJ_ERR     ObjectType = "ERROR"
	OBJ_STR     ObjectType = "STRING"
	OBJ_BUILTIN ObjectType = "BUILTIN"
	OBJ_LIST    ObjectType = "LIST"
)

var (
	TrueBool   = &Boolean{Value: true}
	FalseBool  = &Boolean{Value: false}
	NullObject = &Null{}
)

type Object interface {
	Type() ObjectType
	Inspect() string
}

type Integer struct {
	Value int64
}

func (i *Integer) Inspect() string        { return fmt.Sprintf("%d", i.Value) }
func (i *Integer) Type() ObjectType       { return OBJ_INTEGER }
func NewIntegerObject(val int64) *Integer { return &Integer{Value: val} }

type String struct {
	Value string
}

func (s *String) Inspect() string        { return s.Value }
func (s *String) Type() ObjectType       { return OBJ_STR }
func NewStringObject(val string) *String { return &String{Value: val} }

type Boolean struct {
	Value bool
}

func (b *Boolean) Inspect() string  { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType { return OBJ_BOOLEAN }
func BoolFromNative(val bool) *Boolean {
	if val {
		return TrueBool
	}

	return FalseBool
}

type Function struct {
	Parameters []*ast.Identifier
	Body       *ast.BlockStatement
	Env        *Environment
}

func (f *Function) Inspect() string {
	var out strings.Builder

	ps := []string{}
	for _, p := range f.Parameters {
		ps = append(ps, p.String())
	}

	out.WriteString("fn(")
	out.WriteString(strings.Join(ps, ", "))
	out.WriteString(") {\n")
	out.WriteString(f.Body.String())
	out.WriteString("\n}")

	return out.String()
}
func (f *Function) Type() ObjectType { return OBJ_FUNC }
func NewFunctionObject(
	params []*ast.Identifier,
	body *ast.BlockStatement,
	env *Environment,
) *Function {
	return &Function{
		Parameters: params,
		Body:       body,
		Env:        env,
	}
}

type List struct {
	Elems []Object
}

func (l *List) Inspect() string {
	var out strings.Builder

	es := []string{}
	for _, elem := range l.Elems {
		es = append(es, elem.Inspect())
	}

	out.WriteString("[")
	out.WriteString(strings.Join(es, ", "))
	out.WriteString("]")

	return out.String()
}
func (l *List) Type() ObjectType         { return OBJ_LIST }
func NewListObject(elems []Object) *List { return &List{Elems: elems} }

type Builtin struct {
	Fn BuiltinFn
}

func (b *Builtin) Inspect() string  { return "builtin" }
func (b *Builtin) Type() ObjectType { return OBJ_BUILTIN }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return OBJ_BOOLEAN }

type ReturnVal struct {
	Value Object
}

func (rv *ReturnVal) Inspect() string    { return rv.Value.Inspect() }
func (rv *ReturnVal) Type() ObjectType   { return OBJ_RETURN }
func NewReturnVal(val Object) *ReturnVal { return &ReturnVal{Value: val} }

type Err struct {
	Msg string
}

func (e *Err) Inspect() string  { return "ERROR: " + e.Msg }
func (e *Err) Type() ObjectType { return OBJ_ERR }

func NewErr(format string, args ...interface{}) *Err {
	return &Err{Msg: fmt.Sprintf(format, args...)}
}

func IsErr(obj Object) bool {
	if obj != nil {
		return obj.Type() == OBJ_ERR
	}

	return false
}
