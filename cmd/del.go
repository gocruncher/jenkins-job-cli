package cmd

import (
	"fmt"
	"github.com/ASalimov/jbuilder/cmd/jb"
	"github.com/spf13/cobra"
)

func init(){
	delCmd := &cobra.Command{
		Use:   "del",
		Short: "Delete a particular jenkins settings in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			del(args)
		},
	}
	rootCmd.AddCommand(delCmd)
}

func del(args []string){
	if len(args)==0{
		fmt.Println("which environment should be removed?")
		for {
			name:=getAnswer("name: ")
			if err:=jb.DelEnv(jb.EName(name));err!=nil{
				fmt.Println(err)
				continue
			}
			fmt.Println("Removed.")
			return
		}
	}
	if err:=jb.DelEnv(jb.EName(args[0]));err!=nil{
		fmt.Println(err)
	}
	fmt.Println("Removed.")

}

