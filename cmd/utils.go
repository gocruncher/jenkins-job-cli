package cmd

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
	"strings"
)

func getAnswer(question string, defAnswer string) string {
	for {
		rl, err := readline.New(question)
		defer rl.Close()
		if err != nil {
			panic(err)
		}
		line, err := rl.ReadlineWithDefault(defAnswer)
		line = strings.TrimSpace(line)
		if err != nil {
			fmt.Println(err.Error())
			os.Exit(0)
		}
		if len(line) == 0 { // io.EOF
			continue
		} else {
			return line
		}
	}
}

func findBestChoice(val string, choices []string) string {
	for i := len(val); i > 0; i-- {
		cur := val[:i]
		for _, v := range choices {
			if len(v) >= i && cur == v[:i] {
				return v
			}
		}
	}
	if len(choices) > 0 {
		return choices[0]
	}
	return ""
}

func findBestChoices(val string, choices []string) []string {
	rsp := []string{}
	i := len(val)
	for _, v := range choices {

		if len(v) >= i && val == v[:i] {
			rsp = append(rsp, v)
		}
	}
	return rsp
}
