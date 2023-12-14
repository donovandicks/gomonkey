package object_test

import (
	"testing"

	"github.com/donovandicks/gomonkey/internal/object"
	"github.com/stretchr/testify/assert"
)

func TestObject_Hashable(t *testing.T) {
	t.Parallel()

	cases := []struct {
		name   string
		first  object.HashableObject
		second object.HashableObject
		equal  bool
	}{
		{
			name:   "strings: equal",
			first:  object.NewStringObject("hello"),
			second: object.NewStringObject("hello"),
			equal:  true,
		},
		{
			name:   "strings: unequal",
			first:  object.NewStringObject("hello"),
			second: object.NewStringObject("world"),
			equal:  false,
		},
		{
			name:   "integers: equal",
			first:  object.NewIntegerObject(1),
			second: object.NewIntegerObject(1),
			equal:  true,
		},
		{
			name:   "integers: unequal",
			first:  object.NewIntegerObject(1),
			second: object.NewIntegerObject(2),
			equal:  false,
		},
		{
			name:   "different types",
			first:  object.NewIntegerObject(1),
			second: object.NewStringObject("1"),
			equal:  false,
		},
	}

	for _, testCase := range cases {
		tc := testCase

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			if tc.equal {
				assert.Equal(t, tc.first.Hash(), tc.second.Hash())
			} else {
				assert.NotEqual(t, tc.first.Hash(), tc.second.Hash())
			}
		})
	}
}
