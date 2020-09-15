package jj

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func getEnv(name string) Env {
	for _, e := range GetEnvs() {
		if e.Name == EName(name) {
			return e
		}
	}
	panic("no any envs")
}

func TestInit(t *testing.T) {
	e := getEnv("uat")
	time.Sleep(time.Second * 3)
	GetJobInfo(e, "core-change-zone")
	fmt.Println("jd ", e)
}

func TestGetLastSuccessfulBuildDuration(t *testing.T) {
	rsp, err := GetLastSuccessfulBuildInfo(getEnv("pi"), "config-deploy-manual")
	assert.NoError(t, err)
	fmt.Println("jd ", rsp)
}

func TestCancelJob(t *testing.T) {
	status, err := CancelJob(getEnv("uat"), "web-rpm-build-manual", 40)
	assert.NoError(t, err)
	fmt.Println(status)
}

func TestCancelQueue(t *testing.T) {
	CancelQueue(getEnv("uat"), 657)

}
