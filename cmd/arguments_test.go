package cmd

import (
	"github.com/magiconair/properties/assert"
	"testing"
)

func TestValidate(t *testing.T) {

	type tcase struct {
		arg     string
		isError bool
	}

	tcases := []tcase{
		{"a=b", false},
		{"aaaaaa=", false},
		{"=", true},
		{"==", true},
		{"", true},
		{"aaaa", true},
		{"=bbbb", true},
		{"key", true},
	}
	for _, tc := range tcases {
		t.Run(tc.arg, func(t *testing.T) {
			a := arguments{[]string{tc.arg}}
			err := a.validate()
			if err != nil && !tc.isError {
				t.Error(err)
			}
			if err == nil && tc.isError {
				t.Errorf("there must be an error")
			}
		})
	}
}

func TestGet(t *testing.T) {
	type tcase struct {
		arg  string
		name string
		val  string
	}
	tcases := []tcase{
		{"a=b", "a", "b"},
		{"aaa=", "aaa", ""},
		{"bbb=", "aaa", ""},
	}
	for _, tc := range tcases {
		t.Run(tc.arg, func(t *testing.T) {
			a := arguments{[]string{tc.arg}}
			val, _ := a.get(tc.name)
			assert.Equal(t, tc.val, val)
		})
	}
}
