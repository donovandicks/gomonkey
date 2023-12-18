package opcode_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/opcode"
	"github.com/stretchr/testify/assert"
)

func TestInstruction_New(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name        string
		oc          opcode.OpCode
		operands    []int
		expected    []byte
		expectedStr string
	}{
		{
			name:        "new constant",
			oc:          opcode.OpConstant,
			operands:    []int{65534},
			expected:    []byte{byte(opcode.OpConstant), 255, 254},
			expectedStr: "0000 OpConstant 65534",
		},
	}

	for _, testCase := range cases {
		tc := testCase
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			instr := opcode.NewInstruction(tc.oc, tc.operands)

			assert.Equal(t, tc.expected, instr)
		})
	}
}
