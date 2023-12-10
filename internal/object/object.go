package object

import "fmt"

type ObjectType string

const (
	OBJ_INTEGER ObjectType = "INTEGER"
	OBJ_BOOLEAN ObjectType = "BOOLEAN"
	OBJ_NULL    ObjectType = "NULL"
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

func (b *Boolean) Inspect() string       { return fmt.Sprintf("%t", b.Value) }
func (b *Boolean) Type() ObjectType      { return OBJ_BOOLEAN }
func NewBooleanObject(val bool) *Boolean { return &Boolean{Value: val} }

type Null struct{}

func (n *Null) Inspect() string  { return "null" }
func (n *Null) Type() ObjectType { return OBJ_BOOLEAN }
func NewNullObject() *Null       { return &Null{} }
