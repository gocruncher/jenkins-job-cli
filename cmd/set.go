package cmd

import (
	"fmt"
	"github.com/ASalimov/jbuilder/cmd/jb"
	"github.com/spf13/cobra"
)

func init() {
	setCmd := &cobra.Command{
		Use: `set NAME
where NAME is alias of your jenkins service which will be used in other commands`,
		Short: "Set a jenkins settings in the config file",
		Long: `Set a jenkins settings in the config file. Please make sure that the above command was executed successfully. 
Specifying a NAME that already exists will merge new fields on top of existing values for those fields.
`,
		Run: func(cmd *cobra.Command, args []string) {
			set(args)
		},
		Args: cobra.ExactArgs(1),
	}

	rootCmd.AddCommand(setCmd)
}

func set(args []string) {
	fmt.Println("please specify requires parameters:")
	var name string
	if len(args) == 0 {
		name = getBaseAnswer("name(test,stage, etc.): ", "")

	} else {
		name = args[0]
	}
	_, env := jb.GetEnv(name)

	url := getBaseAnswer("url: ", env.Url)
	var authtype string
	for {
		fmt.Println(`Choose an option from the following list:
	n - No authorization
	a - API token`, authtype)
		authtype = getAnswer("authorization type: ", string(env.Type), []string{"n", "a"})

		if authtype == "n" || authtype == "a" {
			break
		}
	}
	env.Url = url
	env.Name = jb.EName(name)
	env.Type = jb.EType(authtype)
	if authtype == "a" {
		env.Login = getBaseAnswer("login: ", env.Login)
		env.Secret = getBaseAnswer("token: ", env.Secret)

	}
	fmt.Println("checking...")
	err := jb.Check(env)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to get access")
		return
	}
	jb.SetEnv(env)
	fmt.Println("Added")
}
