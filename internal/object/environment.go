package object

type Environment struct {
	vals map[string]Object
}

func NewEnv() *Environment {
	return &Environment{
		vals: make(map[string]Object),
	}
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.vals[name]
	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.vals[name] = val
	return val
}
