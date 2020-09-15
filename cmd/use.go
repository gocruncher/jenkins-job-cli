package cmd

import (
	"fmt"
	"github.com/gocruncher/jenkins-job-cli/cmd/jj"
	"github.com/spf13/cobra"
)

func init() {
	useCmd := &cobra.Command{
		Use:                   "use NAME",
		DisableFlagsInUseLine: true,
		Short:                 "Makes a specific Jenkins name by default",
		Run: func(cmd *cobra.Command, args []string) {
			use(cmd)
		},

		PersistentPreRunE: preRunE,
		Args:              cobra.ExactArgs(1),
	}
	rootCmd.AddCommand(useCmd)
}

func use(cmd *cobra.Command) {
	jj.Init(cmd.Flags().Args()[0])
	jj.SetDef(cmd.Flags().Args()[0])
	fmt.Println(cmd.Flags().Args()[0] + " have been set by default")
}
