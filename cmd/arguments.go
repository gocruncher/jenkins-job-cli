package cmd

import (
	"fmt"
	"strings"
)

type arguments struct {
	args []string
}

func (a arguments) validate() error {
	for _, arg := range a.args {
		if len(arg) < 2 || !strings.Contains(arg, "=") {
			return fmt.Errorf("all argument should look as \"key=val\"")
		}
		if strings.Index(arg, "=") < 1 {
			return fmt.Errorf("all argument should look as \"key=val\"")
		}

	}
	return nil
}

func (a arguments) get(name string) (string, error) {
	for _, arg := range a.args {
		i := strings.Index(arg, name+"=")
		if i > -1 {
			return arg[i+len(name+"="):], nil
		}
	}
	return "", fmt.Errorf("argument is not found")
}
