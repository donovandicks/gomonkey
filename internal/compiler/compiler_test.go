package compiler_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/compiler"
	"github.com/donovandicks/gomonkey/internal/lexer"
	"github.com/donovandicks/gomonkey/internal/opcode"
	"github.com/donovandicks/gomonkey/internal/parser"
	"github.com/stretchr/testify/assert"
)

func TestCompiler(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name           string
		input          string
		expectedInstrs []opcode.Instructions
		expectedConsts []interface{}
	}{
		{
			name:  "add constants",
			input: "1 + 2",
			expectedInstrs: []opcode.Instructions{
				opcode.NewInstruction(opcode.OpConstant, []int{0}),
				opcode.NewInstruction(opcode.OpConstant, []int{1}),
			},
			expectedConsts: []interface{}{1, 2},
		},
	}

	for _, testCase := range cases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			l := lexer.NewLexer(tc.input)
			p := parser.NewParser(l)
			program := p.ParseProgram()

			c := compiler.NewCompiler()
			err := c.Compile(program)

			assert.Nil(t, err, "failed to compile: %v", err)

			b := c.Bytecode()

			for _, exp := range tc.expectedInstrs {
				t.Logf("Expected: %s\t Received: %s\n", exp.String(), b.Instrs.String())
			}

			assert.ElementsMatch(t, tc.expectedConsts, b.Consts, "constants do not match")
			assert.ElementsMatch(t, tc.expectedInstrs, b.Instrs, "instructions do not match")
		})
	}
}
