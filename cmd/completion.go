package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"os"
)

func init(){
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
			runCompletionZsh(os.Stdout,"",rootCmd)
		},
	}
	completionCmd.AddCommand(completionCmdBash)
	completionCmd.AddCommand(completionCmdZsh)
	rootCmd.AddCommand(completionCmd)
}