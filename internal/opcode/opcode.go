package opcode

type OpCode byte

const (
	OpConstant OpCode = iota
)

var (
	opNames = map[OpCode]string{
		OpConstant: "OP_CONST",
	}

	// opWidths is an array with the number of bytes required for each operand
	// corresponding to the index of the array
	opWidths = map[OpCode][]int{
		OpConstant: {2},
	}
)

func (oc OpCode) Name() *string {
	name, ok := opNames[oc]
	if !ok {
		return nil
	}
	return &name
}

func (oc OpCode) OperandWidth() []int {
	width, ok := opWidths[oc]
	if !ok {
		return nil
	}
	return width
}
