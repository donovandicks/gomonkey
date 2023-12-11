package object

func IsTruthy(obj Object) bool {
	switch obj {
	case NullObject:
		return false
	case FalseBool:
		return false
	case TrueBool:
		return true
	default:
		return true
	}
}
