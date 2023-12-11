package object

import "fmt"

type ObjectType string

const (
	OBJ_INTEGER ObjectType = "INTEGER"
	OBJ_BOOLEAN ObjectType = "BOOLEAN"
	OBJ_NULL    ObjectType = "NULL"
	OBJ_RETURN  ObjectType = "RETURN"
	OBJ_ERR     ObjectType = "ERROR"
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
