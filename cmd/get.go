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
	"text/tabwriter"
)



var (
	noheader bool
)

func init() {


	var getCmd = &cobra.Command{
		Use:   "get",
		Short: "Display any resources(settings, views, jobs)",
		Run: func(cmd *cobra.Command, args []string) {
			showAllEnvs()
		},
	}
	getCmd.Flags().BoolVar(&noheader,"no-headers",false,"no-headers")

	for _, env:=range jb.GetEnvs(){
		curEnv:=env
		var envCmd = &cobra.Command{
			Use:   string(env.Name),
			Short: "Get environment info",
			Run: func(cmd *cobra.Command, args []string) {
				showAllJobs(curEnv)
			},
		}
		envCmd.Flags().BoolVar(&noheader,"no-headers",false,"no-headers")


		//for _, view:=range jb.GetBundle(env).Views{
		//	//curView := view
		//	//var viewCmd = &cobra.Command{
		//	//	Use:   string(view.Name),
		//	//	Short: "Get environment info",
		//	//	ValidArgs: []string{},
		//	//	Run: func(cmd *cobra.Command, args []string) {
		//	//		if len(args)>0{
		//	//			os.Exit(1)
		//	//		}
		//	//		showAllJobs(curView)
		//	//	},
		//	//}
		//	//viewCmd.Flags().BoolVar(&noheader,"no-headers",false,"no-headers")
		//	//envCmd.AddCommand(viewCmd)
		//
		//
		//}

		getCmd.AddCommand(envCmd)

	}

	rootCmd.AddCommand(getCmd)
}

func showAllEnvs(){
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if !noheader{
		fmt.Fprintf(w, "%s\t%s\t%s\n","Name", "URL","Authorization")
	}
	for _, e:= range jb.GetEnvs(){
		fmt.Fprintf(w, "%s\t%s\t%s\n",e.Name,e.Url, e.Type)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func showAllView(env jb.Env){
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if !noheader{
		fmt.Fprintf(w, "%s\t%s\n","Name", "URL")
	}
	for _, v:= range jb.GetBundle(env).Views{
		fmt.Fprintf(w, "%s\t%s\n",v.Name,v.URL)
	}
	fmt.Fprintln(w)
	w.Flush()
}

func showAllJobs(env jb.Env){
	w := new(tabwriter.Writer)
	// Format in tab-separated columns with a tab stop of 8.
	w.Init(os.Stdout, 0, 8, 0, '\t', 0)
	if !noheader{
		fmt.Fprintf(w, "%s\t%s\n","Name", "URL")
	}
	for _, view:= range jb.GetBundle(env).Views{
		//fmt.Printf("views: %+v",view)
		for _, j:= range view.Jobs{
			fmt.Fprintf(w, "%s\t%s\n",j.Name,j.URL)
		}
	}
	fmt.Fprintln(w)
	w.Flush()

}