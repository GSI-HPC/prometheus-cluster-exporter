// -*- coding: utf-8 -*-
//
// © Copyright 2023 GSI Helmholtzzentrum für Schwerionenforschung
//
// This software is distributed under
// the terms of the GNU General Public Licence version 3 (GPL Version 3),
// copied verbatim in the file "LICENCE".

package main

import (
	"bytes"
	"io/ioutil"
	"log"
	"os/exec"
	"strings"
	"time"
)

type jobInfo struct {
	jobid   string
	account string
	user    string
}

type runningJobsResult struct {
	elapsed float64
	jobs    []jobInfo
	err     error
}

const SQUEUE = "squeue"

func retrieveRunningJobs(channel chan<- runningJobsResult) {

	start := time.Now()

	_, err := exec.LookPath(SQUEUE)
	if err != nil {
		log.Fatal(err)
	}

	cmd := exec.Command(SQUEUE, "-ah", "-o", "%A %a %u")

	pipe, err := cmd.StdoutPipe()
	if err != nil {
		channel <- runningJobsResult{0, nil, err}
		return
	}

	err = cmd.Start()
	if err != nil {
		channel <- runningJobsResult{0, nil, err}
		return
	}

	out, err := ioutil.ReadAll(pipe)
	if err != nil {
		channel <- runningJobsResult{0, nil, err}
		return
	}

	// TODO Timeout handling?
	err = cmd.Wait()
	if err != nil {
		channel <- runningJobsResult{0, nil, err}
		return
	}

	// TrimSpace on []bytes is more efficient than calling TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(out))

	lines := strings.Split(content, "\n")

	jobs := make([]jobInfo, len(lines))

	for i, line := range lines {
		fields := strings.Fields(line)
		jobs[i] = jobInfo{fields[0], fields[1], fields[2]}
	}

	elapsed := time.Since(start).Seconds()

	channel <- runningJobsResult{elapsed, jobs, nil}
}
