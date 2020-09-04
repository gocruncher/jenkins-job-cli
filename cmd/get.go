/*
Copyright © 2020 NAME HERE <EMAIL ADDRESS>

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/
package cmd

import (
	"fmt"
	"github.com/ASalimov/jbuilder/cmd/jb"
	"github.com/spf13/cobra"
	"os"
	"regexp"
	"strings"
	"text/tabwriter"
	"time"
)

var (
	noheader bool
	ENV      string
)

// -- string Value
//type stringValue string

//func newStringValue(val string, p *string) *stringValue {
//	*p = val
//	return (*stringValue)(p)
//}
//
//func (s *stringValue) Set(val string) error {
//	*s = stringValue(val)
//	return nil
//}
//func (s *stringValue) Type() string {
//	return "string"
//}
//func (s *stringValue) String() string { return string(*s) }

func init() {

	annotation := make(map[string][]string)
	annotation[cobra.BashCompCustom] = []string{"__jb_get_env"}

	//flag := &pflag.Flag{
	//	Name:        "env",
	//	Usage:       "Current environment",
	//	Annotations: annotation,
	//	Value:newStringValue("", &ENV),
	//}
	var getCmd = &cobra.Command{
		Use:   "get [names|jobs]",
		Short: "Display any resources(settings, jobs)",
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 && ENV != "" {
				showAllJobs(jb.Init(ENV))
				os.Exit(0)
			}
			if len(args) > 1 && args[0] == "compline" {
				space := regexp.MustCompile(`\s+`)
				s := space.ReplaceAllString(args[1], " ")
				s = strings.ReplaceAll(s, "=", " ")
				s = strings.ReplaceAll(s, "--name", "-n")
				if i := strings.Index(s, "-n"); i > 0 {
					s1 := strings.Split(s[i+3:], " ")
					showJobs(strings.TrimSpace(s1[0]))
				} else {
					showJobs(string(jb.GetDefEnv()))
				}
				os.Exit(0)
			}
			showAllEnvs()
			os.Exit(0)
			return nil
		},
		PreRunE: preRunE,
	}

	//getCmd.Flags().AddFlag(flag)
	getCmd.Flags().BoolVar(&noheader, "no-headers", false, "no-headers")
	getCmd.Flags().StringVarP(&ENV, "name", "n", "", "current Jenkins name")
	//getCmd.Flags().StringVarP(&ENV, "env", "e", "", "")
	//for _, env:=range jb.GetEnvs(){
	//	curEnv:=env
	//	var envCmd = &cobra.Command{
	//		Use:   string(env.Name),
	//		Short: "Get environment info",
	//		Run: func(cmd *cobra.Command, args []string) {
	//			showAllJobs(curEnv)
	//		},
	//	}
	//	envCmd.Flags().BoolVar(&noheader,"no-headers",false,"no-headers")
	//
	//	getCmd.AddCommand(envCmd)
	//
	//}

	rootCmd.AddCommand(getCmd)
}

func showJobs(eName string) {
	ch := make(chan struct{}, 1)

	// Run your long running function in it's own goroutine and pass back it's
	// response into our channel.
	go func() {
		env := jb.Init(eName)
		showAllJobs(env)
		ch <- struct{}{}
	}()

	select {
	case <-ch:
		os.Exit(0)
	case <-time.After(100 * time.Millisecond):
		os.Exit(0)
	}
}

func showAllEnvs() {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if !noheader {
		fmt.Fprintf(w, "%s\t%s\t%s\n", "Name", "URL", "Authorization")
	}
	for _, e := range jb.GetEnvs() {
		fmt.Fprintf(w, "%s\t%s\t%s\n", e.Name, e.Url, e.Type)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func showAllJobs(env jb.Env) {
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if !noheader {
		fmt.Fprintf(w, "%s\t%s\n", "Name", "URL")
	}
	for _, view := range jb.GetBundle(env).Views {
		//fmt.Printf("views: %+v",view)
		for _, j := range view.Jobs {
			fmt.Fprintf(w, "%s\t%s\n", j.Name, j.URL)
		}
	}
	fmt.Fprintln(w)
	w.Flush()

}
