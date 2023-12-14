package interpreter

import (
	"fmt"

	"github.com/donovandicks/gomonkey/internal/object"
)

var Builtins = map[string]*object.Builtin{
	"len": {
		Fn: Len,
	},
	"print": {
		Fn: Print,
	},
}

func Len(args ...object.Object) object.Object {
	if len(args) != 1 {
		return object.NewErr("invalid number of args %d, expected 1", len(args))
	}

	switch arg := args[0].(type) {
	case *object.String:
		return object.NewIntegerObject(int64(len(arg.Value)))
	case *object.List:
		return object.NewIntegerObject(int64(len(arg.Elems)))
	default:
		return object.NewErr("invalid argument %s", args[0].Type())
	}
}

func Print(args ...object.Object) object.Object {
	for _, arg := range args {
		fmt.Println(arg.Inspect())
	}

	return nil
}
