package cmd

import (
	"fmt"
	"github.com/gocruncher/jenkins-job-cli/cmd/jj"
	"github.com/spf13/cobra"
)

var URL string
var login string
var token string

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
	setCmd.Flags().StringVarP(&URL, "url", "u", "", "URL of the Jenkins")
	setCmd.Flags().StringVarP(&login, "login", "l", "", "login")
	setCmd.Flags().StringVarP(&token, "token", "t", "", "API token")

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
	_, env := jj.GetEnv(name)
	if URL == "" {
		env.Url = getBaseAnswer("url: ", env.Url)
	} else {
		env.Url = URL

	}

	var authtype string
	if URL != "" {
		authtype = "a"
	} else {
		for {
			fmt.Println(`Choose an option from the following list:
	n - No authorization
	a - API token`, authtype)
			authtype = getAnswer("authorization type: ", string(env.Type), []string{"n", "a"})

			if authtype == "n" || authtype == "a" {
				break
			}
		}
	}
	env.Name = jj.EName(name)
	env.Type = jj.EType(authtype)
	if authtype == "a" {
		if login == "" {
			env.Login = getBaseAnswer("login: ", env.Login)
		} else {
			env.Login = login

		}
		if token == "" {
			env.Secret = getBaseAnswer("token: ", env.Secret)
		} else {
			env.Secret = token
		}
	}
	fmt.Println("checking...")
	err := jj.Check(env)
	if err != nil {
		fmt.Println(err.Error())
		fmt.Println("Failed to get access")
		return
	}
	jj.SetEnv(env)
	fmt.Println("Added")
}
