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
	"errors"
	"io/ioutil"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

const GETENT = "getent"

type userInfo struct {
	user string
	uid  int
	gid  int
}

type groupInfo struct {
	group string
	gid   int
}

type userInfoMap map[int]userInfo
type groupInfoMap map[int]groupInfo

type userInfoMapResult struct {
	elapsed float64
	users   userInfoMap
	err     error
}

type groupInfoMapResult struct {
	elapsed float64
	groups  groupInfoMap
	err     error
}

func createUserInfoMap(channel chan<- userInfoMapResult) {

	start := time.Now()

	userInfoMap := make(userInfoMap)

	cmd := exec.Command(GETENT, "passwd")

	pipe, err := cmd.StdoutPipe()

	if err != nil {
		channel <- userInfoMapResult{0, nil, err}
		return
	}

	if err := cmd.Start(); err != nil {
		channel <- userInfoMapResult{0, nil, err}
		return
	}

	out, err := ioutil.ReadAll(pipe)

	if err != nil {
		channel <- userInfoMapResult{0, nil, err}
		return
	}

	// TODO Timeout handling?
	if err := cmd.Wait(); err != nil {
		channel <- userInfoMapResult{0, nil, err}
		return
	}

	// TrimSpace on []bytes is more efficient than calling TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(out))

	if len(content) == 0 {
		channel <- userInfoMapResult{0, nil, errors.New("retrieved content in createUserInfoMap() is empty")}
		return
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {

		fields := strings.SplitN(line, ":", 5)

		if len(fields) < 4 {
			channel <- userInfoMapResult{0, nil, errors.New("insufficient field count found in line: " + line)}
			return
		}

		user := fields[0]

		uid, err := strconv.Atoi(fields[2])
		if err != nil {
			channel <- userInfoMapResult{0, nil, err}
			return
		}

		gid, err := strconv.Atoi(fields[3])
		if err != nil {
			channel <- userInfoMapResult{0, nil, err}
			return
		}

		userInfoMap[uid] = userInfo{user, uid, gid}
	}

	elapsed := time.Since(start).Seconds()

	channel <- userInfoMapResult{elapsed, userInfoMap, nil}
}

func createGroupInfoMap(channel chan<- groupInfoMapResult) {

	start := time.Now()

	groupInfoMap := make(groupInfoMap)

	cmd := exec.Command(GETENT, "group")

	pipe, err := cmd.StdoutPipe()

	if err != nil {
		channel <- groupInfoMapResult{0, nil, err}
		return
	}

	if err := cmd.Start(); err != nil {
		channel <- groupInfoMapResult{0, nil, err}
		return
	}

	out, err := ioutil.ReadAll(pipe)

	if err != nil {
		channel <- groupInfoMapResult{0, nil, err}
		return
	}

	// TODO Timeout handling?
	if err := cmd.Wait(); err != nil {
		channel <- groupInfoMapResult{0, nil, err}
		return
	}

	// TrimSpace on []bytes is more efficient than calling TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(out))

	if len(content) == 0 {
		channel <- groupInfoMapResult{0, nil, errors.New("retrieved content in createGroupInfoMap() is empty")}
		return
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		fields := strings.SplitN(line, ":", 4)

		if len(fields) < 3 {
			channel <- groupInfoMapResult{0, nil, errors.New("insufficient field count found in line: " + line)}
			return
		}

		group := fields[0]

		gid, err := strconv.Atoi(fields[2])
		if err != nil {
			channel <- groupInfoMapResult{0, nil, err}
			return
		}

		groupInfoMap[gid] = groupInfo{group, gid}
	}

	elapsed := time.Since(start).Seconds()

	channel <- groupInfoMapResult{elapsed, groupInfoMap, nil}
}
