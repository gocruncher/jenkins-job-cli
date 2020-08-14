package cmd

import (
	"fmt"
	"github.com/ASalimov/jbuilder/cmd/jb"
	"github.com/spf13/cobra"
)


func init() {
	setCmd := &cobra.Command{
		Use:   "set",
		Short: "Set a jenkins settings in the config file",
		Run: func(cmd *cobra.Command, args []string) {
			set(args)
		},
	}
	rootCmd.AddCommand(setCmd)
}

func set(args []string) {
	fmt.Println("please specify requires parameters:")
	name:= getName(args)
	url:=getAnswer("url: ")
	var authtype string
	for {
		fmt.Println(`Choose an option from the following list:
	n - No authorization
	a - API token`,authtype)
		authtype=getAnswer("authorization type: ")

		if authtype =="n" ||authtype=="a"{
			break
		}
	}
	env := jb.Env{Url: url,Name: jb.EName(name),Type: jb.EType(authtype)}
	if authtype=="a"{
		env.Login =getAnswer("login: ")
		env.Secret =getAnswer("token: ")

	}
	fmt.Println("checking...")
	err:=jb.Check(env)
	if err!=nil{
		fmt.Println(err.Error())
		fmt.Println("Failed to get access")
		return
	}
	jb.SetEnv(env)
}
