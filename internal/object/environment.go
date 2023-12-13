package object

type Environment struct {
	vals  map[string]Object
	outer *Environment
}

func NewEnv() *Environment {
	return &Environment{
		vals: make(map[string]Object),
	}
}

func NewEnvFromEnv(outer *Environment) *Environment {
	return &Environment{
		vals:  make(map[string]Object),
		outer: outer,
	}
}

func (e *Environment) Values() map[string]Object {
	return e.vals
}

func (e *Environment) With(vals map[string]Object) *Environment {
	for ident, val := range vals {
		e.Set(ident, val)
	}

	return e
}

func (e *Environment) Get(name string) (Object, bool) {
	val, ok := e.vals[name]
	if !ok && e.outer != nil {
		return e.outer.Get(name)
	}

	return val, ok
}

func (e *Environment) Set(name string, val Object) Object {
	e.vals[name] = val
	return val
}

func (e *Environment) Update(name string, val Object) Object {
	_, ok := e.vals[name]
	if !ok {
		return NewErr("undefined variable '%s'", name)
	}

	e.vals[name] = val
	return val
}
