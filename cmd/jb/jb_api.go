package jb

import (
	"encoding/json"
	"errors"
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strconv"
	"strings"
)

// External API

func GetBundle(env Env) *Bundle {
	for _, b := range bundles {
		if b.Name == env.Name {
			return b
		}
	}
	return nil
}

func DelEnv(name EName) error {
	for i, e := range config.Envs {
		if e.Name == name {
			config.Envs = append(config.Envs[:i], config.Envs[i+1:]...)
			SetEnv(e)
			return nil
		}
	}
	return fmt.Errorf("'%s' name is not found", name)
}

func Check(env Env) error {
	code, rspbin, _, err := req(env, "api/json", []byte{})
	if err != nil {
		return err
	}
	if code != 200 {
		return fmt.Errorf("http code: %d\nresponse: %s", code, rspbin)
	}
	return nil
}

func GetEnvs() []Env {
	return config.Envs
}

func GetEnv(name string) (error, Env) {
	var env Env
	if name == "" {
		name = string(GetDefEnv())
	}
	for _, e := range GetEnvs() {
		if e.Name == EName(name) {
			env = e
			break
		}
	}
	if env == (Env{}) {
		return ErrNoEnv, env
	}
	return nil, env
}

func (ji *JobInfo) GetParameterDefinitions() []ParameterDefinitions {
	for _, j := range ji.Property {
		if len(j.ParameterDefinitions) > 0 {
			return j.ParameterDefinitions
		}
	}
	return []ParameterDefinitions{}
}

func GetDefEnv() EName {
	if config.Use == "" {
		return GetEnvs()[0].Name
	}
	return config.Use
}
func SetDef(eName string) {
	var env Env
	for _, e := range GetEnvs() {
		if e.Name == EName(eName) {
			env = e
			break
		}
	}
	if env == (Env{}) {
		panic("Environment " + eName + " is not found or could not be initialised")
	}
	config.Use = env.Name
	SetConf()
}

func SetConf() {
	out, _ := yaml.Marshal(config)
	if _, err := os.Stat(homeDir); os.IsNotExist(err) {
		err := os.MkdirAll(homeDir, os.ModePerm)
		if err != nil {
			panic(err)
		}
	}
	err := ioutil.WriteFile(homeDir+configFile, out, 0644)
	if err != nil {
		panic(err)
	}
}

func SetEnv(env Env) {
	for i, e := range config.Envs {
		if e.Name == env.Name {
			config.Envs[i] = env
			break
		}
	}
	config.Envs = append(config.Envs, env)
	SetConf()

}

//func GetJobInfo(env Env, jobName string) JobInfo{
//	var jobinfo JobInfo
//	code, rsp, _, err := req(env,"job/"+jobName+"/api/json", []byte{})
//	if err != nil {
//		panic(err)
//	}
//	if code != 200 {
//		panic("failed to get job details,code" + strconv.Itoa(code) + ", " + string(rsp))
//	}
//	err = json.Unmarshal(rsp, &jobinfo)
//	if err!=nil{
//		panic("failed to get Job information")
//	}
//	return jobinfo
//}
func GetJobInfo(env Env, jobName string) (error, *JobInfo) {
	bundle := GetBundle(env)
	var jobInfo *JobInfo
	for _, ji := range bundle.JobsInfo {
		if ji.Name == jobName {
			jobInfo = &ji
			break
		}
	}
	var fetchJobInfo = func() (error, JobInfo) {
		var ji JobInfo
		code, rsp, _, err := req(env, "job/"+jobName+"/api/json", []byte{})
		if err != nil {
			panic(err)
		}
		if code != 200 {
			return ErrNoJob, ji
		}
		err = json.Unmarshal(rsp, &ji)
		if err != nil {
			panic("failed to get Job information")
		}
		mutex.Lock()
		defer mutex.Unlock()
		bundle.JobsInfo = append(bundle.JobsInfo, ji)
		updateCache(env, bundle)
		return nil, ji
	}
	if jobInfo != nil {
		//fmt.Println("async")
		go fetchJobInfo()
	} else {
		//fmt.Println("sync")
		err, ji := fetchJobInfo()
		if err != nil {
			return err, jobInfo
		}
		jobInfo = &ji
	}
	return nil, jobInfo

}

//
//func GetJobParameterDefinitions(env Env, jobName string) []ParameterDefinitions {
//	bundle:=GetBundle(env)
//	params:= []ParameterDefinitions{}
//	for _,jp:= range bundle.JobsParameters{
//		if jp.Name == jobName{
//			params = jp.Parameters
//		}
//	}
//
//	var fetchParameters = func() []ParameterDefinitions{
//		jobinfo:=GetJobInfo(env,jobName)
//		parameters:=jobinfo.GetParameterDefinitions()
//		mutex.Lock()
//		defer mutex.Unlock()
//		bundle.JobsParameters = append(bundle.JobsParameters,JobsParameters{Name: jobName, Parameters: parameters})
//		updateCache(env, bundle)
//		return parameters
//	}
//
//	if len(params)>0{
//		fmt.Println("async")
//		go fetchParameters()
//	}else{
//		fmt.Println("sync")
//		return fetchParameters()
//	}
//	return params
//
//
//
//}
//

func GetBuildInfo(env Env, job string, id int) (*BuildInfo, error) {
	code, rsp, _, err := req(env, "job/"+job+"/"+strconv.Itoa(id)+"/api/json", []byte{})
	if err != nil {
		panic(err)
	}
	if code != 200 {
		return nil, errors.New("failed to get job details,code" + strconv.Itoa(code) + ", " + string(rsp))
	}
	var bi BuildInfo
	err = json.Unmarshal(rsp, &bi)
	if err != nil {
		return nil, err
	}
	return &bi, nil
}

func GetLastSuccessfulBuildInfo(env Env, job string) (*BuildInfo, error) {
	code, rsp, _, err := req(env, "job/"+job+"/lastSuccessfulBuild/api/json", []byte{})
	if err != nil {
		panic(err)
	}
	if code != 200 {
		return nil, errors.New("failed to get job details,code" + strconv.Itoa(code) + ", " + string(rsp))
	}
	var bi BuildInfo
	err = json.Unmarshal(rsp, &bi)
	if err != nil {
		return nil, err
	}
	return &bi, nil
}

func Build(env Env, job string, query string) string {
	target := "/build"
	if len(query) > 0 {
		target = "/buildWithParameters?" + query
	}
	code, rsp, headers, err := req(env, "job/"+job+target, []byte{})
	if err != nil {
		panic(err)
	}
	if code != 201 {
		panic(errors.New("failed to start job details,code" + strconv.Itoa(code) + ", " + string(rsp)))
	}
	location := headers["Location"][0]
	splitedUrl := strings.Split(location, "/")
	return splitedUrl[len(splitedUrl)-2]

}

func CancelQueue(env Env, id int) {
	req(env, "queue/cancelItem?id="+strconv.Itoa(id), []byte{})
}
func CancelJob(env Env, job string, id int) (string, error) {
	code, _, _, err := req(env, "job/"+job+"/"+strconv.Itoa(id)+"/stop", []byte{})
	if err != nil {
		panic(err)
	}
	if code != 200 {
		return "", errors.New("failed to cancel the job,code" + strconv.Itoa(code))
	}
	bi, err := GetBuildInfo(env, job, id)
	if err != nil {
		return "", err
	}
	return bi.Result, nil
}

func Console(env Env, job string, id int, start string) (string, string, error) {
	//web-rpm-build-manual/149/logText/progressiveHtml
	code, rsp, h, err := req(env, "job/"+job+"/"+strconv.Itoa(id)+"/logText/progressiveHtml", []byte("start="+start))
	if err != nil {
		return "", "", err
	}
	if code != 200 {
		return "", "", errors.New("code = " + strconv.Itoa(code))
	}
	//fmt.Println(h)
	return string(rsp), h["X-Text-Size"][0], nil
}

func GetQueueInfo(env Env, id int) (error, QueueInfo) {
	var queueInfo QueueInfo
	code, rsp, _, err := req(env, "/queue/item/"+strconv.Itoa(id)+"/api/json", []byte{})
	if err != nil {
		panic(err)
	}
	if code != 200 {
		return errors.New("failed to get queue details,code" + strconv.Itoa(code) + ", " + string(rsp)), QueueInfo{}
	}
	err = json.Unmarshal(rsp, &queueInfo)
	if err != nil {
		return errors.New("failed to get Job information"), QueueInfo{}
	}
	return nil, queueInfo
}

func GetQueues(env Env) Queues {
	var queues Queues
	code, rsp, _, err := req(env, "/queue/api/json", []byte{})
	if err != nil {
		panic(err)
	}
	if code != 200 {
		panic("failed to get queue list,code" + strconv.Itoa(code) + ", " + string(rsp))
	}
	err = json.Unmarshal(rsp, &queues)
	if err != nil {
		panic("failed to get Queue list")
	}
	return queues
}
