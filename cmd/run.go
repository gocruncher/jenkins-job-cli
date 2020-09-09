package cmd

import (
	"errors"
	"fmt"
	"github.com/ASalimov/bar"
	"github.com/ASalimov/jbuilder/cmd/jb"
	"github.com/chzyer/readline"
	"github.com/spf13/cobra"
	"github.com/ttacon/chalk"
	"net/url"
	"os"
	"os/signal"
	"strconv"
	"strings"
	"sync"
	"time"
)

var usageTamplate = `Usage:{{if .Runnable}}
  {{.UseLine}}{{end}}{{if .HasAvailableSubCommands}}
  {{.CommandPath}} [command]{{end}}{{if gt (len .Aliases) 0}}

Aliases:
  {{.NameAndAliases}}{{end}}{{if .HasExample}}

Examples:
{{.Example}}{{end}}{{if .HasAvailableSubCommands}}

Available Commands:{{range .Commands}}{{if (or .IsAvailableCommand (eq .Name "help"))}}
  {{rpad .Name .NamePadding }} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableLocalFlags}}

Flags:
{{.LocalFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasAvailableInheritedFlags}}

Global Flags:
{{.InheritedFlags.FlagUsages | trimTrailingWhitespaces}}{{end}}{{if .HasHelpSubCommands}}

Additional help topics:{{range .Commands}}{{if .IsAdditionalHelpTopicCommand}}
  {{rpad .CommandPath .CommandPathPadding}} {{.Short}}{{end}}{{end}}{{end}}{{if .HasAvailableSubCommands}}

Use "{{.CommandPath}} [command] --help" for more information about a command.{{end}}
`

var listenerStatus bool

type st struct {
	name  string
	id    int
	queue int
}

var curSt st
var barMutex sync.Mutex
var closeCh chan struct{}
var stdinListener *jbStdin

func init() {
	var runCmd = &cobra.Command{
		Use:     "run JOB",
		Aliases: []string{"r"},
		Short:   "Run the specified jenkins job",
		Run: func(cmd *cobra.Command, args []string) {
			run_job(args[0])
		},

		Args:         cobra.ExactArgs(1),
		PreRunE:      preRunE,
		SilenceUsage: false,
	}
	runCmd.Flags().StringVarP(&ENV, "name", "n", "", "current Jenkins name")
	runCmd.SetUsageTemplate(usageTamplate)
	rootCmd.AddCommand(runCmd)
}

func run_job(name string) {
	env := jb.Init(ENV)
	fmt.Printf("Job will be started in %s environment\n", env.Name)
	fmt.Println(env.Url + "/job/" + name)

	bar.InitTerminal()
	data := map[string]string{}
	err, jobInfo := jb.GetJobInfo(env, name)
	if err != nil {
		if err == jb.ErrNoJob {
			fmt.Printf("Error: job '%s' does not exist", name)
			os.Exit(1)
		}
		panic(err)
	}
	params := jobInfo.GetParameterDefinitions()
	if len(params) == 0 {
		rl, err := readline.New("Press any key to continue: ")
		defer rl.Close()
		_, err = rl.Readline()
		if err != nil {
			panic(err)
		}
	}
	for _, pd := range params {
		cline := ""
		defVal := pd.DefaultParameterValue.Value
		curChoices := pd.Choices
		for {
			rl, err := NewReadLine(chalk.Underline.TextStyle(pd.Name)+": ", pd.Choices)
			defer rl.Close()
			if err != nil {
				panic(err)
			}
			if pd.Type == "ChoiceParameterDefinition" {
				defVal = ""
			}
			line, err := rl.ReadlineWithDefault(defVal)
			line = strings.TrimSpace(line)
			if err != nil { // io.EOF
				os.Exit(1)
			}
			if pd.Type == "ChoiceParameterDefinition" {

				for _, val := range pd.Choices {
					if line == val {
						cline = val
						break
					}
				}
				if cline == "" {
					curChoices = findBestChoices(line, pd.Choices)
					if len(curChoices) == 0 {
						curChoices = pd.Choices
					} else if len(curChoices) == 1 {
						defVal = curChoices[0]
					} else {
						defVal = line
					}
					for _, val := range curChoices {
						fmt.Printf("%s\t", val)
					}
					if len(curChoices) > 0 {
						fmt.Println()
					}

					continue
				}
			} else {
				cline = line
			}
			break
		}

		data[pd.Name] = cline
	}

	urlquery := url.Values{}
	for key, val := range data {
		urlquery.Add(key, val)
	}
	queueId := jb.Build(env, name, urlquery.Encode())

	keyCh := make(chan string)
	stdinListener = NewStdin()
	go listenKeys(keyCh)
	go listenInterrupt(env)
	queueId1, _ := strconv.Atoi(queueId)
	curSt.queue = queueId1
	curSt.name = name
	number := waitForExecutor(env, queueId1)
	curSt.id = number
	err = watchTheJob(env, name, number, keyCh)
	if err != nil {
		return
	}
	curSt = st{}
	for _, jChild := range jobInfo.DownstreamProjects {
		err = watchNext(env, name, jChild.Name, number, keyCh)
		if err != nil {
			return
		}
		curSt = st{}
	}
	fmt.Println(chalk.Green.Color("done"))
	return
}

func waitForExecutor(env jb.Env, queueId int) int {
	informed := false
	for {
		err, queueInfo := jb.GetQueueInfo(env, queueId)
		if err != nil {
			panic(err)
		}
		if !queueInfo.Blocked && queueInfo.Executable.URL != "" {
			return queueInfo.Executable.Number
		} else {
			if !informed {
				//clearer := strings.Repeat(" ", int(110)-1)
				fmt.Println("waiting for next available executor..  ")
				informed = true
			}
			time.Sleep(100 * time.Millisecond)
		}
	}
}

func barHandler(jobUrl string, keyCh chan string, chMsg chan string, finishCh chan struct {
	err    error
	result string
}, wg *sync.WaitGroup) {
	defer wg.Done()
	barMutex.Lock()
	fmt.Print("\033[F")
	br := bar.NewWithOpts(
		bar.WithDimensions(100, 20),
		bar.WithLines(1),
		bar.WithFormat(
			fmt.Sprintf(
				"%srunning...%s :percent :bar %s:eta%s",
				chalk.White,
				chalk.Reset,
				chalk.Green,
				chalk.Reset)))
	br.Tick()
	barMutex.Unlock()
	for {
		select {
		case stdin, _ := <-keyCh:
			if []byte(stdin)[0] == 10 {
				barMutex.Lock()
				br.SetLines(br.GetLines() + 1)
				barMutex.Unlock()
			}
		case msg := <-chMsg:

			if msg != "" {
				barMutex.Lock()
				br.Interrupt(msg)
				barMutex.Unlock()
			} else {
				barMutex.Lock()
				br.Tick()
				barMutex.Unlock()
			}

		case info := <-finishCh:
			if info.err != nil && br.GetLines() < 5 {
				for {
					if br.GetLines() < 10 {
						barMutex.Lock()
						br.SetLines(br.GetLines() + 1)
						barMutex.Unlock()
					} else {
						break
					}
				}

			}
			if info.err == nil {
				fmt.Printf("\r%s", strings.Repeat(" ", int(50)-1))
				fmt.Print("\033[F")
			}

			barMutex.Lock()
			br.SetFormat(fmt.Sprintf(jobUrl + ": " + info.result))
			br.Done()
			barMutex.Unlock()
			if info.err != nil {
				fmt.Println(chalk.Red.Color("failed"))
			}
			return
		case <-closeCh:
			return
		}

	}
}

func watchTheJob(env jb.Env, name string, number int, keyCh chan string) error {
	jobUrl := env.Url + "/job/" + name + "/" + strconv.Itoa(number) + "/console"
	lastBuild, _ := jb.GetLastSuccessfulBuildInfo(env, name)
	listenerStatus = true
	defer func() {
		listenerStatus = false
	}()
	ticks := 1
	t := 0
	cursor := "0"
	stime := getTime()
	chMsg := make(chan string)
	closeCh = make(chan struct{})
	finishCh := make(chan struct {
		err    error
		result string
	})
	var wg sync.WaitGroup
	wg.Add(1)
	go barHandler(jobUrl, keyCh, chMsg, finishCh, &wg)
	defer close(closeCh)
	defer wg.Wait()
	handle := func(cursor string) string {
		output, nextCursor, err := jb.Console(env, name, number, cursor)
		if err != nil || cursor == nextCursor {
			return cursor
		}
		j := 1
		for {
			lines := strings.Split(output, "\n")

			count := len(lines)
			if count > 50 {
				count = 50
			}
			for i := count; i >= 1; i-- {
				rline := []rune(string(lines[len(lines)-i]))
				if err != nil { // io.EOF
					break
				}
				j := 0
				size := 100
				for {
					var fline string
					s := j * size
					e := (j + 1) * size
					if len(rline) > e {
						fline = string(rline[s:e])
					} else {
						fline = string(rline[s:len(rline)])
					}
					if len(strings.TrimSpace(fline)) > 0 {
						chMsg <- fline
						time.After(40 * time.Millisecond)

					}
					j++
					if len(rline) <= e || len(rline) > 10*size {
						break
					}
				}
			}
			break
			j++
			//time.Sleep(4*time.Millisecond)
		}
		return nextCursor
	}
	for t > -1 {

		if t%5 == 0 && t > 1 {
			curBuild, err := jb.GetBuildInfo(env, name, number)
			if err != nil {
				if getTime()-stime > int64(30*time.Millisecond) {
					err := errors.New("failed")
					finishCh <- struct {
						err    error
						result string
					}{err, err.Error()}
					return err
				}
			} else {
				if !curBuild.Building {
					if curBuild.Result == "SUCCESS" {
						for {
							nc := handle(cursor)
							if nc == cursor {
								break
							}
							cursor = nc
						}
						finishCh <- struct {
							err    error
							result string
						}{nil, curBuild.Result}
						return nil
					} else {
						err := errors.New("failed")
						finishCh <- struct {
							err    error
							result string
						}{err, curBuild.Result}
						return err
					}
				}
			}
		}

		ncursor := handle(cursor)
		if ncursor != cursor {
			cursor = ncursor
			ctime := getTime()
			dtime := ctime - stime
			newTicks := int(float64(dtime) / float64(lastBuild.Duration) * 100)
			for ticks < newTicks && ticks < 99 {
				chMsg <- ""
				ticks++
			}
			if dtime < 500 {
				time.Sleep(time.Duration(500-dtime) * time.Millisecond)
			}
		} else {
			time.Sleep(100 * time.Millisecond)
		}
		t++
	}
	err := errors.New("failed")
	finishCh <- struct {
		err    error
		result string
	}{err, "timeout"}
	return err
}

func watchNext(env jb.Env, parentName string, childName string, parentJobID int, keyCh chan string) error {
	for i := 0; ; i++ {
		bi, err := findDownstreamInBuilds(env, parentName, childName, parentJobID)
		if err != nil {
			queueId, err := findDownstreamInQueue(env, parentName, childName, parentJobID)
			curSt.queue = queueId
			curSt.name = childName
			if err != nil {
				time.Sleep(250 * time.Millisecond)
				continue
			}
			number := waitForExecutor(env, queueId)
			curSt.id = number
			return watchTheJob(env, childName, number, keyCh)
		} else {
			id, _ := strconv.Atoi(bi.Id)
			curSt.name = childName
			curSt.id = id
			return watchTheJob(env, childName, id, keyCh)
		}
	}

}

func findDownstreamInBuilds(env jb.Env, parentName string, childName string, parent int) (*jb.BuildInfo, error) {
	err, jobInfo := jb.GetJobInfo(env, childName)
	if err != nil {
		panic(err)
	}
	number := jobInfo.LastBuild.Number
	for i := 5; i >= 0; i-- {
		bi, err := jb.GetBuildInfo(env, childName, number-i)
		if err != nil {
			continue
		}
		for _, a := range bi.Actions {
			for _, c := range a.Causes {
				if c.UpstreamBuild == parent && c.UpstreamProject == parentName {
					return bi, nil
				}
			}
		}
	}
	return &jb.BuildInfo{}, errors.New("not found")
}

func findDownstreamInQueue(env jb.Env, parentName string, childName string, parentJobID int) (int, error) {
	queues := jb.GetQueues(env)
	for _, queue := range queues.Items {
		if queue.Task.Name == childName {
			for _, action := range queue.Actions {
				for _, cause := range action.Causes {
					if cause.UpstreamBuild == parentJobID && cause.UpstreamProject == parentName {
						return queue.ID, nil
					}
				}
			}
		}
	}
	return 0, errors.New("not found")
}

func listenKeys(out chan string) {
	stdinListener.NewListener()
	bt := make([]byte, 1)
	for {
		n, err := stdinListener.Read(bt)
		if err != nil || n == 0 {
			return
		}
		barMutex.Lock()
		if listenerStatus {
			out <- string(bt)
		}
		barMutex.Unlock()
	}

}

func listenInterrupt(env jb.Env) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt)
	go func() {
		for _ = range c {
			if curSt.name != "" {
				barMutex.Lock()
				stdinListener.NewListener()
				readline.Stdin = stdinListener
				rl, err := readline.New(fmt.Sprintf("There is active build: %s. Do you want to cancel it [Y/n]:", curSt.name))
				defer rl.Close()
				if err != nil {
					panic(err)
				}
				line, err := rl.Readline()
				if err != nil { // io.EOF
					os.Exit(1)
				}
				if line == "Y" || line == "y" {

					if curSt.queue != 0 {
						fmt.Println("canceling queue...")
						jb.CancelQueue(env, curSt.queue)
					}
					if curSt.id != 0 {
						fmt.Println("canceling job...")
						status, err := jb.CancelJob(env, curSt.name, curSt.id)
						if err != nil {
							fmt.Printf("failed to cancel job, error %s", err)
							os.Exit(0)
						}
						if status != "ABORTED" {
							fmt.Printf("Job already has been executed, status: %s", status)
							os.Exit(0)
						}
						fmt.Println("Canceled")
						os.Exit(0)
					}
					if curSt.queue != 0 && curSt.id == 0 {
						err, jobInfo := jb.GetJobInfo(env, curSt.name)
						if err != nil {
							panic(err)
						}
						number := jobInfo.LastBuild.Number
						for i := 0; i < 3; i++ {
							bi, err := jb.GetBuildInfo(env, curSt.name, number-i)
							if err != nil {
								continue
							}
							if bi.QueueId == curSt.queue {
								if bi.Result != "ABORTED" {
									fmt.Printf("Job already has been executed, status: %s", bi.Result)
									os.Exit(0)
								} else {
									fmt.Println("Canceled!")
									os.Exit(0)
								}
							}
						}
						fmt.Println("Canceled!!!")
					}

				}
				os.Exit(0)

			}

		}
	}()

}
