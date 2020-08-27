/*
Copyright © 2020 asalimov

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
	"bytes"
	"fmt"
	"github.com/spf13/cobra"
	"io"
	"os"
	"time"
)

const (
	bash_completion_func = `__jb_parse_get()
{
    local jb_output out
	echo "will be run jb get --no-headers ${nouns[@]} 2>/dev/null" >> /tmp/jblogs
    if jb_output=$(jb get --no-headers "${nouns[@]}" 2>/dev/null); then
        out=($(echo "${jb_output}" | awk '{print $1}'))
        COMPREPLY=( $( compgen -W "${out[*]}" -- "$cur" ) )
    fi
}

__jb_get_env()
{
    return 0
}

__jb_get_resource()
{
    # if [[ ${#nouns[@]} -eq 0 ]]; then
	#	echo "step3" >> /tmp/jblogs
	#	echo "step3 ${#nouns[@]}" >> /tmp/jblogs
    #    return 1
    # fi
	echo "step4 ${#nouns[@]} ${nouns[@]}" >> /tmp/jblogs
	__jb_parse_get ${nouns[@]}
    if [[ $? -eq 0 ]]; then
        return 0
    fi
}

__jb_custom_func() {
	echo "step0 ${last_command} and ${flags_with_completion[@]}|${nouns[@]}|${cur}" >> /tmp/jblogs
    case ${last_command} in
        jb_get | jb_run | jb_delete | jb_stop)
			echo "step2" >> /tmp/jblogs
            __jb_get_resource
            return
            ;;
        *)
            ;;
    esac
}
`
)

var cfgFile string

const defaultBoilerPlate = `
# Copyright 2016 The Kubernetes Authors.
#
# Licensed under the Apache License, Version 2.0 (the "License");
# you may not use this file except in compliance with the License.
# You may obtain a copy of the License at
#
#     http://www.apache.org/licenses/LICENSE-2.0
#
# Unless required by applicable law or agreed to in writing, software
# distributed under the License is distributed on an "AS IS" BASIS,
# WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
# See the License for the specific language governing permissions and
# limitations under the License.
`

// rootCmd represents the base command when called without any subcommands
var rootCmd = &cobra.Command{
	Use:   "jb",
	Short: "jb - simple command line utility which just runs any jenkins job",
	Long: `jb is simple command line utility which just runs any jenkins job.
Before you start, please configure access to the Jenkins service using "jb set" command. For more information you can type "jb help set"
After that, you can enable shell autocompletion for convenient work. To do this, run following:
 # For zsh completion	
   echo 'source <(jb completion zsh)' >>~/.zshrc
 # For bash completion
   echo 'source <(jb completion bash)' >>~/.bashrc
  
`,
	//ValidArgs: []string{"run","get","set","del","completion"},
	BashCompletionFunction: bash_completion_func,

	// Uncomment the following line if your bare application
	// has an action associated with it:
	//	Run: func(cmd *cobra.Command, args []string) { },
	Example: `	jb run app-build
	jb run -e prod web-build`,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}

func getTime() int64 {
	return time.Now().UnixNano() / int64(time.Millisecond)
}

func runCompletionZsh(out io.Writer, boilerPlate string, root *cobra.Command) error {
	zshHead := "#compdef jb\n"

	out.Write([]byte(zshHead))

	if len(boilerPlate) == 0 {
		boilerPlate = defaultBoilerPlate
	}
	if _, err := out.Write([]byte(boilerPlate)); err != nil {
		return err
	}

	zshInitialization := `
__jb_bash_source() {
	alias shopt=':'
	emulate -L sh
	setopt kshglob noshglob braceexpand
	source "$@"
}
__jb_type() {
	# -t is not supported by zsh
	if [ "$1" == "-t" ]; then
		shift
		# fake Bash 4 to disable "complete -o nospace". Instead
		# "compopt +-o nospace" is used in the code to toggle trailing
		# spaces. We don't support that, but leave trailing spaces on
		# all the time
		if [ "$1" = "__jb_compopt" ]; then
			echo builtin
			return 0
		fi
	fi
	type "$@"
}
__jb_compgen() {
	local completions w
	completions=( $(compgen "$@") ) || return $?
	# filter by given word as prefix
	while [[ "$1" = -* && "$1" != -- ]]; do
		shift
		shift
	done
	if [[ "$1" == -- ]]; then
		shift
	fi
	for w in "${completions[@]}"; do
		if [[ "${w}" = "$1"* ]]; then
			echo "${w}"
		fi
	done
}
__jb_compopt() {
	true # don't do anything. Not supported by bashcompinit in zsh
}
__jb_ltrim_colon_completions()
{
	if [[ "$1" == *:* && "$COMP_WORDBREAKS" == *:* ]]; then
		# Remove colon-word prefix from COMPREPLY items
		local colon_word=${1%${1##*:}}
		local i=${#COMPREPLY[*]}
		while [[ $((--i)) -ge 0 ]]; do
			COMPREPLY[$i]=${COMPREPLY[$i]#"$colon_word"}
		done
	fi
}
__jb_get_comp_words_by_ref() {
	cur="${COMP_WORDS[COMP_CWORD]}"
	prev="${COMP_WORDS[${COMP_CWORD}-1]}"
	words=("${COMP_WORDS[@]}")
	cword=("${COMP_CWORD[@]}")
}
__jb_filedir() {
	# Don't need to do anything here.
	# Otherwise we will get trailing space without "compopt -o nospace"
	true
}
autoload -U +X bashcompinit && bashcompinit
# use word boundary patterns for BSD or GNU sed
LWORD='[[:<:]]'
RWORD='[[:>:]]'
if sed --help 2>&1 | grep -q 'GNU\|BusyBox'; then
	LWORD='\<'
	RWORD='\>'
fi
__jb_convert_bash_to_zsh() {
	sed \
	-e 's/declare -F/whence -w/' \
	-e 's/_get_comp_words_by_ref "\$@"/_get_comp_words_by_ref "\$*"/' \
	-e 's/local \([a-zA-Z0-9_]*\)=/local \1; \1=/' \
	-e 's/flags+=("\(--.*\)=")/flags+=("\1"); two_word_flags+=("\1")/' \
	-e 's/must_have_one_flag+=("\(--.*\)=")/must_have_one_flag+=("\1")/' \
	-e "s/${LWORD}_filedir${RWORD}/__jb_filedir/g" \
	-e "s/${LWORD}_get_comp_words_by_ref${RWORD}/__jb_get_comp_words_by_ref/g" \
	-e "s/${LWORD}__ltrim_colon_completions${RWORD}/__jb_ltrim_colon_completions/g" \
	-e "s/${LWORD}compgen${RWORD}/__jb_compgen/g" \
	-e "s/${LWORD}compopt${RWORD}/__jb_compopt/g" \
	-e "s/${LWORD}declare${RWORD}/builtin declare/g" \
	-e "s/\\\$(type${RWORD}/\$(__jb_type/g" \
	<<'BASH_COMPLETION_EOF'
`
	out.Write([]byte(zshInitialization))

	buf := new(bytes.Buffer)
	root.GenBashCompletion(buf)
	out.Write(buf.Bytes())

	zshTail := `
BASH_COMPLETION_EOF
}
__jb_bash_source <(__jb_convert_bash_to_zsh)
`
	out.Write([]byte(zshTail))
	return nil
}
