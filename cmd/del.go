package cmd

import (
	"fmt"
	"github.com/gocruncher/jenkins-job-ctl/cmd/jj"
	"github.com/spf13/cobra"
)

func init() {
	delCmd := &cobra.Command{
		Use:   "del NAME",
		Short: "Delete a particular jenkins settings in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			del(args)
		},
	}
	rootCmd.AddCommand(delCmd)
}

func del(args []string) {
	if len(args) == 0 {
		fmt.Println("which Jenkins should be removed?")
		choices := []string{}
		for _, e := range jj.GetEnvs() {
			choices = append(choices, string(e.Name))
		}
		for {

			name := getAnswer("name: ", "", choices)
			if err := jj.DelEnv(jj.EName(name)); err != nil {
				fmt.Println(err)
				continue
			}
			fmt.Println("Removed.")
			return
		}
	}
	if err := jj.DelEnv(jj.EName(args[0])); err != nil {
		fmt.Println(err)
	} else {
		fmt.Println("Removed.")
	}

}
