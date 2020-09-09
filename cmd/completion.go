package cmd

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

func init() {
	// completionCmd represents the completion command
	var completionCmd = &cobra.Command{
		Use:   "completion SHELL",
		Short: "Create a bash/zsh completion script",
		Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hey ")
		},
	}
	var completionCmdBash = &cobra.Command{
		Use:   "bash",
		Short: "Generates bash completion scripts",
		Long: `To load completion run

. <(bitbucket completion)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			rootCmd.GenBashCompletion(os.Stdout)
		},
	}
	var completionCmdZsh = &cobra.Command{
		Use:   "zsh",
		Short: "Generates zsh completion scripts",
		Long: `To load completion run

. <(bitbucket completion zsh)

To configure your bash shell to load completions for each session add to your bashrc

# ~/.bashrc or ~/.profile
. <(bitbucket completion)
`,
		Run: func(cmd *cobra.Command, args []string) {
			runCompletionZsh(os.Stdout, "", rootCmd)
		},
	}
	var completionCmdCheck = &cobra.Command{
		Use:   "check",
		Short: "Checks you OS completions settings",
		Long:  `Checks you OS completions settings`,
		Run: func(cmd *cobra.Command, args []string) {
			check()
		},
	}
	completionCmd.AddCommand(completionCmdBash)
	completionCmd.AddCommand(completionCmdZsh)
	completionCmd.AddCommand(completionCmdCheck)
	rootCmd.AddCommand(completionCmd)
}

func check() {
	shell := os.Getenv("SHELL")
	if strings.Contains(shell, "zsh") {
		checkZsh()
		return
	}
	if strings.Contains(shell, "bash") {
		checkBash(shell)
		return
	}
	fmt.Printf("shell '%s' is not supported yet", shell)

}
func checkBash(shell string) {
	fmt.Println("==> current shell is bash - OK")
	cmd := exec.Command(shell, "--version")
	verInfoB, err := cmd.Output()
	if err != nil {
		panic(err)
	}
	verInfo := string(verInfoB)
	i := strings.Index(verInfo, "version ") + 8
	ver := verInfo[i : i+3]
	verA := strings.Split(verInfo[i:i+3], ".")
	if verA[0] > "4" || (verA[0] == "4" && verA[1] >= "1") {
		fmt.Printf("==> bash version is %s - OK\n", ver)
	} else {
		fmt.Printf("==> bash version is %s - FAILED\n", ver)
		fmt.Printf(
			`
Error: completion script doesn't work correctly with Bash %s. It requires bash-completion v2 and Bash 4.1+. 
Thus, to be able to correctly use completion on macOS, you have to install and use Bash 4.1+. The following instruction helps you with it:
https://itnext.io/upgrading-bash-on-macos-7138bd1066ba
`, ver)
		return
	}
	fmt.Println("Checked. After reloading your shell, jb autocompletion should be working.")
}

func checkZsh() {
	homeDir, _ := os.UserHomeDir()
	fmt.Println("==> current shell is zsh - OK")
	zshrc, err := ioutil.ReadFile(homeDir + "/.zshrc")
	scanner := bufio.NewScanner(bytes.NewReader(zshrc))
	var buffer strings.Builder
	for scanner.Scan() {
		text := scanner.Text()
		trimText := strings.TrimSpace(text)
		if len(trimText) > 0 && trimText[:1] != "#" {
			buffer.WriteString(text)
			buffer.WriteString("\n")
		}
	}
	if err != nil {
		fmt.Println(
			`
Error: ~/.zshrc file was not found. As the next step, please create file and add the following line:
  source <(jb completion zsh)`)
		return
	}
	fmt.Println("==> file ~/.zshrc is exist - OK")
	if !strings.Contains(buffer.String(), "jb completion zsh") {
		fmt.Println("==> completion script has not been added - FAILED")
		fmt.Println(
			`
Error: ~/.zshrc should have line with jb completion script for the shell. Make sure, that the following line has been added the the ~/.zshrc file:
    source <(jb completion zsh)`)
		return
	}
	fmt.Println("==> completion script has been added - OK")
	cmd := exec.Command("zsh", "-c", "source ~/.zshrc; compdef")
	stderr, err := cmd.StderrPipe()
	log.SetOutput(os.Stderr)
	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}
	rsp, _ := ioutil.ReadAll(stderr)
	if strings.Contains(string(rsp), "not found") {
		compinitPos := strings.Index(buffer.String(), "autoload -Uz compinit")
		scriptPos := strings.Index(buffer.String(), "jb completion zsh")
		if compinitPos > -1 && compinitPos > scriptPos {
			fmt.Println("==> the initialized compdef script was added incorrectly - FAILED")
			fmt.Println(
				`
Error: the initialized compdef script was added incorrectly. Move the following to the beginning of your ~/.zshrc file
    autoload -Uz compinit
    compinit`)
			return
		} else {
			fmt.Println("==> compdef functions is not exist - FAILED")
			fmt.Println(
				`
Error: current shell does not have required 'compdef' function. Add the following to the beginning of your ~/.zshrc file:
    autoload -Uz compinit
    compinit`)
			return
		}

	}
	fmt.Println("==> compdef functions is exist - OK")

	fmt.Println("Checked. After reloading your shell, jb autocompletion should be working.")
}
