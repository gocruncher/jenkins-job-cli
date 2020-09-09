package cmd

import (
	"fmt"
	"github.com/chzyer/readline"
	"os"
	"strings"
)

func filterInput(r rune) (rune, bool) {
	switch r {
	// block CtrlZ feature
	case readline.CharCtrlZ:
		return 3, true
	}
	return r, true
}
func getBaseAnswer(question string, defAnswer string) string {
	return getAnswer(question, defAnswer, []string{})
}
func NewReadLine(question string, choices []string) (*readline.Instance, error) {
	completer := []readline.PrefixCompleterInterface{}
	for _, choice := range choices {
		completer = append(completer, readline.PcItem(choice))
	}
	rl, err := readline.NewEx(&readline.Config{
		Prompt:              question,
		HistoryFile:         "/tmp/readline.tmp",
		AutoComplete:        readline.NewPrefixCompleter(completer...),
		InterruptPrompt:     "^C",
		EOFPrompt:           "exit",
		HistorySearchFold:   true,
		FuncFilterInputRune: filterInput,
	})
	return rl, err
}
func getAnswer(question string, defAnswer string, choices []string) string {

	for {
		rl, err := NewReadLine(question, choices)
		defer rl.Close()

		if err != nil {
			panic(err)
		}
		line, err := rl.ReadlineWithDefault(defAnswer)
		fmt.Println("entered: ", line)
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
