package jj

import (
	"crypto/tls"
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"strconv"
	"strings"
	"sync"
	"time"
)

var homeDir string
var Debug = false

const cacheFile = "cache"
const configFile = "config.yaml"

var config Config
var bundles []*Bundle
var mutex sync.Mutex
var ErrNoEnv = errors.New("no env")
var ErrNoJob = errors.New("no job")

func init() {
	homeDir, _ = os.UserHomeDir()
	homeDir = homeDir + "/.jj/"
	initConfig()
}

func Init(envName string) Env {
	err, env := GetEnv(envName)
	if err != nil {
		panic(ErrNoEnv)
	}
	initBundle(env)
	return env

}

type EName string
type EType string

func (t EType) String() string {
	if t == EType("n") {
		return "none"
	}
	return "undefined"
}

type Config struct {
	Use  EName `yaml:"use"`
	Envs []Env `yaml:"envs"`
}

type Bundle struct {
	Name     EName     `json:"name"`
	Views    []View    `json:"views"`
	JobsInfo []JobInfo `json:"JobsInfo"`
}

type Job struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type View struct {
	Name string `json:"name"`
	URL  string `json:"url"`
	Desc string `json:"desc"`
	Jobs []Job  `json:"jobs"`
}

type BuildInfo struct {
	Id      string `json:"id"`
	Actions []struct {
		Parameters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters,omitempty"`
		Causes []struct {
			ShortDescription string `json:"shortDescription"`
			UpstreamBuild    int    `json:"upstreamBuild"`
			UpstreamProject  string `json:"upstreamProject"`
			UpstreamURL      string `json:"upstreamUrl"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
	Duration int    `json:"duration"`
	Building bool   `json:"building"`
	Result   string `json:"result"`
	QueueId  int    `json:"queueId"`
}

type ParameterDefinitions struct {
	DefaultParameterValue struct {
		Name  string `json:"name"`
		Value string `json:"value"`
	} `json:"defaultParameterValue"`
	Description string   `json:"description"`
	Name        string   `json:"name"`
	Type        string   `json:"type"`
	Choices     []string `json:"choices,omitempty"`
}

type Env struct {
	Url    string `yaml:"url"`
	Name   EName  `yaml:"name"`
	Type   EType  `yaml:"type"`
	Login  string `yaml:"login"`
	Secret string `yaml:"secret"`
}

type JobInfo struct {
	Name               string `json:"name"`
	NextBuildNumber    int    `json:"nextBuildNumber"`
	DownstreamProjects []struct {
		Name string `json:"name"`
		URL  string `json:"url"`
	} `json:"downstreamProjects"`
	LastBuild struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"lastBuild"`
	LastCompletedBuild struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	}
	InQueue  bool `json:"inQueue"`
	Property []struct {
		ParameterDefinitions []ParameterDefinitions `json:"parameterDefinitions,omitempty"`
	} `json:"property"`
}

type QueueInfo struct {
	Actions []struct {
		Parameters []struct {
			Name  string `json:"name"`
			Value string `json:"value"`
		} `json:"parameters,omitempty"`
		Causes []struct {
			ShortDescription string      `json:"shortDescription"`
			UserID           interface{} `json:"userId"`
			UserName         string      `json:"userName"`
		} `json:"causes,omitempty"`
	} `json:"actions"`
	Blocked      bool   `json:"blocked"`
	Buildable    bool   `json:"buildable"`
	ID           int    `json:"id"`
	InQueueSince int64  `json:"inQueueSince"`
	Params       string `json:"params"`
	Stuck        bool   `json:"stuck"`
	Task         struct {
		Name  string `json:"name"`
		URL   string `json:"url"`
		Color string `json:"color"`
	} `json:"task"`
	URL        string      `json:"url"`
	Why        interface{} `json:"why"`
	Cancelled  bool        `json:"cancelled"`
	Executable struct {
		Number int    `json:"number"`
		URL    string `json:"url"`
	} `json:"executable"`
}

type Queues struct {
	DiscoverableItems []interface{} `json:"discoverableItems"`
	Items             []struct {
		Actions []struct {
			Parameters []struct {
				Name  string `json:"name"`
				Value string `json:"value"`
			} `json:"parameters,omitempty"`
			Causes []struct {
				ShortDescription string `json:"shortDescription"`
				UpstreamBuild    int    `json:"upstreamBuild"`
				UpstreamProject  string `json:"upstreamProject"`
				UpstreamURL      string `json:"upstreamUrl"`
			} `json:"causes,omitempty"`
		} `json:"actions"`
		Blocked                    bool   `json:"blocked"`
		Buildable                  bool   `json:"buildable"`
		BuildableStartMilliseconds int64  `json:"buildableStartMilliseconds"`
		ID                         int    `json:"id"`
		InQueueSince               int64  `json:"inQueueSince"`
		Params                     string `json:"params"`
		Stuck                      bool   `json:"stuck"`
		Task                       struct {
			Color string `json:"color"`
			Name  string `json:"name"`
			URL   string `json:"url"`
		} `json:"task"`
		URL string `json:"url"`
		Why string `json:"why"`
	} `json:"items"`
}

func initConfig() {
	data, err := ioutil.ReadFile(homeDir + configFile)
	if err != nil {
		return
	}
	err = yaml.Unmarshal(data, &config)
	changed := false
	for {
	Repeat:
		for i, env1 := range config.Envs {
			for j, env2 := range config.Envs {
				if env1.Name == env2.Name && i != j {
					config.Envs[len(config.Envs)-1], config.Envs[j] = config.Envs[j], config.Envs[len(config.Envs)-1]
					config.Envs = config.Envs[:len(config.Envs)-1]
					changed = true
					break Repeat
				}
			}
		}
		if changed {
			SetConf()
			changed = false
		} else {
			break
		}
	}
	if err != nil {
		log.Fatalf("error: %v", err)
	}
}

func setViews(env Env, views []View) {
	mutex.Lock()
	defer mutex.Unlock()
	bundle := GetBundle(env)
	if bundle == nil {
		bundle = &Bundle{Name: env.Name, Views: views, JobsInfo: []JobInfo{}}
		bundles = append(bundles, bundle)
	} else {
		bundle.Views = views
	}
	updateCache(env, bundle)
}

func initBundle(env Env) {
	var bundle Bundle
	cachebin, err := ioutil.ReadFile(homeDir + cacheFile + "." + string(env.Name))
	err = json.Unmarshal(cachebin, &bundle)
	if err != nil {
		fetchBundle(env, false)
	} else {
		mutex.Lock()
		defer mutex.Unlock()
		bundles = append(bundles, &bundle)
		go fetchBundle(env, true)
	}
}

func fetchBundle(env Env, async bool) {
	code, rspbin, _, err := req(env, "POST", "api/json", []byte{})
	if err != nil {
		if !async {
			fmt.Printf("error: %s\n", err)
		}
		return
	}
	if code != 200 {
		if !async {
			fmt.Println("failed to get job details:code " + strconv.Itoa(code) + " " + env.Login)
		}

		return
	}
	var rsp struct {
		Views []View `json:"views"`
	}
	err = json.Unmarshal(rspbin, &rsp)
	if err != nil {
		panic(err)
	}
	for i, view := range rsp.Views {
		code, vrspbin, _, err := req(env, "POST", "view/"+view.Name+"/api/json", []byte{})
		if err != nil {
			panic(err)
		}
		if code != 200 {
			panic("failed to get view details,code" + strconv.Itoa(code) + ", " + string(rspbin))
		}
		err = json.Unmarshal(vrspbin, &rsp.Views[i])
		if err != nil {
			panic(err)
		}
	}
	setViews(env, rsp.Views)

}

func updateCache(env Env, bundle *Bundle) {
	cachebin, err := json.Marshal(bundle)
	if err == nil {
		path := homeDir + cacheFile + "." + string(env.Name)
		ioutil.WriteFile(path, cachebin, 0644)
	} else {
		fmt.Println("failed to Marshal cache info, ", err.Error())
	}
}
func reqPOST(env Env, method, path string, body []byte) (int, []byte, map[string][]string, error) {
	return req(env, method, path, body)
}
func req(env Env, method, path string, body []byte) (int, []byte, map[string][]string, error) {
	base_url := env.Url
	if base_url[len(base_url)-1:] != "/" {
		base_url = base_url + "/"
	}

	url := base_url + path
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	client := &http.Client{Transport: tr, Timeout: time.Second * 30}
	request, err := http.NewRequest(method, url, strings.NewReader(string(body)))
	if err != nil {
		return 0, nil, nil, err
	}
	request.Header.Add("Accept-Language", "en-us")
	request.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	if env.Type == "a" {
		request.SetBasicAuth(env.Login, env.Secret)
	}
	response, err := client.Do(request)
	if err != nil {
		return 0, nil, nil, err
	}
	defer response.Body.Close()
	contents, err := ioutil.ReadAll(response.Body)
	if err != nil {
		return 0, nil, nil, err
	}

	if Debug {
		fmt.Println("req: ", url, "\t data: ", string(body), "type:", response.Header.Get("Content-Type"), "\n rsp: ", string(contents))
	}
	return response.StatusCode, contents, response.Header, nil
}
