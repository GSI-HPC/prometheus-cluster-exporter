// Copyright 2020 Gabriele Iannetti <g.iannetti@gsi.de>
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
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"strings"
)

type jobInfo struct {
	jobid   string
	account string
	user    string
}

const squeueBin = "/usr/bin/squeue"

func retrieveRunningJobs() ([]jobInfo, error) {

	if _, err := os.Stat(squeueBin); os.IsNotExist(err) {
		log.Fatal(err)
	}

	cmd := exec.Command(squeueBin, "-ah", "-o", "%A %a %u")

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

	lines := strings.Split(content, "\n")

	jobs := make([]jobInfo, len(lines))

	for i, line := range lines {
		fields := strings.Fields(line)
		jobs[i] = jobInfo{fields[0], fields[1], fields[2]}
	}

	return jobs, nil
}
