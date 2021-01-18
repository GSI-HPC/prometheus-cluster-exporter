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

	log "github.com/sirupsen/logrus"
)

const getentBin = "/usr/bin/getent"

type UserInfo struct {
	user string
	uid  int
	gid  int
}

type GroupInfo struct {
	group string
	gid   int
}

type UserInfoMap map[int]UserInfo

type GroupInfoMap map[int]GroupInfo

func createUserInfoMap() (UserInfoMap, error) {

	var m UserInfoMap
	m = make(UserInfoMap)

	if _, err := os.Stat(getentBin); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cmd := exec.Command(getentBin, "passwd")

	pipe, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	out, err := ioutil.ReadAll(pipe)

	// TODO Timeout handling?
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	// TrimSpace on []bytes is more efficient than calling TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(out))

	if len(content) == 0 {
		return nil, errors.New("Retrieved content in createUserInfoMap() is empty")
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {

		fields := strings.SplitN(line, ":", 5)

		if len(fields) < 4 {
			return nil, errors.New("Insufficient field count found in line: " + line)
		}

		user := fields[0]

		uid, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}

		gid, err := strconv.Atoi(fields[3])
		if err != nil {
			return nil, err
		}

		m[uid] = UserInfo{user, uid, gid}
	}

	return m, nil
}

func createGroupInfoMap() (GroupInfoMap, error) {

	var m GroupInfoMap
	m = make(GroupInfoMap)

	if _, err := os.Stat(getentBin); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cmd := exec.Command(getentBin, "group")

	pipe, err := cmd.StdoutPipe()

	if err != nil {
		return nil, err
	}

	if err := cmd.Start(); err != nil {
		return nil, err
	}

	out, err := ioutil.ReadAll(pipe)

	// TODO Timeout handling?
	if err := cmd.Wait(); err != nil {
		return nil, err
	}

	// TrimSpace on []bytes is more efficient than calling TrimSpace on a string since it creates a copy
	content := string(bytes.TrimSpace(out))

	if len(content) == 0 {
		return nil, errors.New("Retrieved content in createGroupInfoMap() is empty")
	}

	lines := strings.Split(content, "\n")

	for _, line := range lines {
		fields := strings.SplitN(line, ":", 4)

		if len(fields) < 3 {
			return nil, errors.New("Insufficient field count found in line: " + line)
		}

		group := fields[0]

		gid, err := strconv.Atoi(fields[2])
		if err != nil {
			return nil, err
		}

		m[gid] = GroupInfo{group, gid}
	}

	return m, nil
}
