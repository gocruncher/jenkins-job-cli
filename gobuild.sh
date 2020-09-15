#!/bin/bash

set -eux

export GOPATH="$(pwd)/.gobuild"
SRCDIR="${GOPATH}/src/github.com/gocruncher/jenkins-job-ctl"

[ -d ${GOPATH} ] && rm -rf ${GOPATH}
mkdir -p ${GOPATH}/{src,pkg,bin}
mkdir -p ${SRCDIR}
cp main.go ${SRCDIR}
cp -r ./cmd ${SRCDIR}
(
    echo ${GOPATH}
    cd ${SRCDIR}
    go get -d ./...
    go install .
)