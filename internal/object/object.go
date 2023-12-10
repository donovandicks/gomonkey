package object

import "fmt"

type ObjectType string

const (
	OBJ_INTEGER ObjectType = "INTEGER"
	OBJ_BOOLEAN ObjectType = "BOOLEAN"
	OBJ_NULL    ObjectType = "NULL"
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
