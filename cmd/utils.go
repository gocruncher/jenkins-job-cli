package cmd

import (
	"github.com/chzyer/readline"
	"strings"
)

func getName(args []string) string{
	var name string
	if len(args)==0{
		name = getAnswer("name(test,stage, etc.): ")

	}else{
		name = args[0]
	}
	return name
}

func getAnswer(question string, ) string{
	for{
		rl, err := readline.New( question)
		if err != nil {
			panic(err)
		}
		defer rl.Close()
		line, err := rl.Readline()
		line = strings.TrimSpace(line)
		if err != nil || len(line)==0 { // io.EOF
			continue
		}else{
			return line
		}
	}
}

