package opcode

import (
	"encoding/binary"
	"fmt"
	"strings"
)

type Instructions []byte

func NewInstruction(code OpCode, operands []int) []byte {
	opWidths := code.OperandWidth()
	if code.Name() == nil || opWidths == nil {
		return []byte{}
	}

	instrLen := 1
	for _, width := range opWidths {
		instrLen += width
	}

	instr := make([]byte, instrLen)
	instr[0] = byte(code)

	offset := 1
	for i, op := range operands {
		width := opWidths[i]
		switch width {
		case 2:
			binary.BigEndian.PutUint16(instr[offset:], uint16(op))
		}
		offset += width
	}

	return instr
}

func (instr Instructions) String() string {
	var out strings.Builder

	i := 0
	for i < len(instr) {
		opCode := OpCode(instr[i])
		opWidths := opCode.OperandWidth()
		opName := opCode.Name()
		if opWidths == nil || opName == nil {
			fmt.Fprintf(&out, "ERROR: undefined opcode %b", instr[i])
			continue
		}

		operands, read := instr.Operands(opWidths)
		fmt.Fprintf(&out, "%04d %s\n", i, instr.Format(*opName, opWidths, operands))

		i += 1 + read // incr by the opcode and the subsequent bytes read
	}

	return out.String()
}

func (instr Instructions) Format(opName string, opWidths []int, operands []int) string {
	opCount := len(opWidths)

	if len(opWidths) != opCount {
		return fmt.Sprintf("ERROR: opcount %d does not match number of operands %d", opCount, len(operands))
	}

	switch opCount {
	case 1:
		return fmt.Sprintf("%s %d", opName, operands[0])
	default:
		return fmt.Sprintf("ERROR: unimplemented for %d operands", opCount)
	}
}

func (instr Instructions) Operands(opWidths []int) ([]int, int) {
	operands := make([]int, len(opWidths))
	offset := 0

	for i, width := range opWidths {
		switch width {
		case 2:
			operands[i] = int(binary.BigEndian.Uint16(instr[offset:]))
		}

		offset += width
	}

	return operands, offset
}
