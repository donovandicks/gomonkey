package compiler

import (
	"github.com/donovandicks/gomonkey/internal/ast"
	"github.com/donovandicks/gomonkey/internal/object"
	"github.com/donovandicks/gomonkey/internal/opcode"
)

type Compiler struct {
	instrs opcode.Instructions
	consts []object.Object
}

func NewCompiler() *Compiler {
	return &Compiler{
		instrs: opcode.Instructions{},
		consts: []object.Object{},
	}
}

func (c *Compiler) Compile(node ast.Node) error {
	return nil
}

func (c *Compiler) Bytecode() *Bytecode {
	return &Bytecode{
		Instrs: c.instrs,
		Consts: c.consts,
	}
}
