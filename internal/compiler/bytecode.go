package compiler

import (
	"github.com/donovandicks/gomonkey/internal/object"
	"github.com/donovandicks/gomonkey/internal/opcode"
)

type Bytecode struct {
	Instrs opcode.Instructions
	Consts []object.Object
}
