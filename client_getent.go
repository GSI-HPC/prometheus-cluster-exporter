// Copyright 2021 Gabriele Iannetti <g.iannetti@gsi.de>
//
// This program is free software: you can redistribute it and/or modify
// it under the terms of the GNU General Public License as published by
// the Free Software Foundation, either version 3 of the License, or
// (at your option) any later version.
//
// This program is distributed in the hope that it will be useful,
// but WITHOUT ANY WARRANTY; without even the implied warranty of
// MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
// GNU General Public License for more details.
//
// You should have received a copy of the GNU General Public License
// along with this program.  If not, see <http://www.gnu.org/licenses/>.

package main

import (
	"bytes"
	"errors"
	"io/ioutil"
	"os"
	"os/exec"
	"strconv"
	"strings"
	"time"

	log "github.com/sirupsen/logrus"
)

const getentBin = "/usr/bin/getent"

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

	if _, err := os.Stat(getentBin); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cmd := exec.Command(getentBin, "passwd")

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

	if _, err := os.Stat(getentBin); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cmd := exec.Command(getentBin, "group")

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
