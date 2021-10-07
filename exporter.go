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
	"errors"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/buger/jsonparser"
	"github.com/prometheus/client_golang/prometheus"
	log "github.com/sirupsen/logrus"
)

type exporter struct {
	channelRunningJobs           chan runningJobsResult
	channelUserInfo              chan userInfoMapResult
	channelGroupInfo             chan groupInfoMapResult
	scrapeActive                 bool
	scrapeMutex                  sync.Mutex
	requestTimeout               int
	urlLustreMetadataOperations  string
	urlLustreJobReadBytes        string
	urlLustreJobWriteBytes       string
	scrapeOKMetric               prometheus.Gauge
	stageExecutionMetric         *prometheus.GaugeVec
	jobMetadataOperationsMetric  *prometheus.GaugeVec
	jobReadThroughputMetric      *prometheus.GaugeVec
	jobWriteThroughputMetric     *prometheus.GaugeVec
	procMetadataOperationsMetric *prometheus.GaugeVec
	procReadThroughputMetric     *prometheus.GaugeVec
	procWriteThroughputMetric    *prometheus.GaugeVec
}

type metadataInfo struct {
	jobid      string
	operations int64
}

type throughputInfo struct {
	jobid      string
	throughput float64
}

func newGaugeVecMetric(namespace string, metricName string, docString string, constLabels []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: namespace,
			Name:      metricName,
			Help:      docString,
		},
		constLabels,
	)
}

func newExporter(requestTimeout int, urlLustreMetadataOperations string, urlLustreJobReadBytes string, urlLustreJobWriteBytes string) *exporter {

	channelRunningJobs := make(chan runningJobsResult)
	channelUserInfo := make(chan userInfoMapResult)
	channelGroupInfo := make(chan groupInfoMapResult)

	if requestTimeout <= 0 {
		log.Fatal("Request timeout must be greater then 0")
	}

	scrapeOKMetric := prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace: namespaceInternals,
		Name:      "scrape_ok",
		Help:      "Indicates if the scrape of the exporter was successful or not.",
	})

	stageExecutionMetric := newGaugeVecMetric(
		namespaceInternals,
		"stage_execution_seconds",
		"Execution duration in seconds spend in a specific exporter stage.",
		[]string{"name"})

	jobMetadataOperationsMetric := newGaugeVecMetric(
		namespace,
		"job_metadata_operations",
		"Total metadata operations of all jobs per account and user.",
		[]string{"account", "user"})

	jobReadThroughputMetric := newGaugeVecMetric(
		namespace,
		"job_read_throughput_bytes",
		"Total IO read throughput of all jobs per account and user in bytes per second.",
		[]string{"account", "user"})

	jobWriteThroughputMetric := newGaugeVecMetric(
		namespace,
		"job_write_throughput_bytes",
		"Total IO write throughput of all jobs per account and user in bytes per second.",
		[]string{"account", "user"})

	procMetadataOperationsMetric := newGaugeVecMetric(
		namespace,
		"proc_metadata_operations",
		"Total metadata operations of process names per group and user.",
		[]string{"proc_name", "group_name", "user_name"})

	procReadThroughputMetric := newGaugeVecMetric(
		namespace,
		"proc_read_throughput_bytes",
		"Total IO read throughput of process names per group and user in bytes per second.",
		[]string{"proc_name", "group_name", "user_name"})

	procWriteThroughputMetric := newGaugeVecMetric(
		namespace,
		"proc_write_throughput_bytes",
		"Total IO write throughput of process names per group and user in bytes per second.",
		[]string{"proc_name", "group_name", "user_name"})

	return &exporter{
		channelRunningJobs:           channelRunningJobs,
		channelUserInfo:              channelUserInfo,
		channelGroupInfo:             channelGroupInfo,
		requestTimeout:               requestTimeout,
		urlLustreMetadataOperations:  urlLustreMetadataOperations,
		urlLustreJobReadBytes:        urlLustreJobReadBytes,
		urlLustreJobWriteBytes:       urlLustreJobWriteBytes,
		scrapeOKMetric:               scrapeOKMetric,
		stageExecutionMetric:         stageExecutionMetric,
		jobMetadataOperationsMetric:  jobMetadataOperationsMetric,
		jobReadThroughputMetric:      jobReadThroughputMetric,
		jobWriteThroughputMetric:     jobWriteThroughputMetric,
		procMetadataOperationsMetric: procMetadataOperationsMetric,
		procReadThroughputMetric:     procReadThroughputMetric,
		procWriteThroughputMetric:    procWriteThroughputMetric,
	}
}

func (e *exporter) Collect(ch chan<- prometheus.Metric) {

	scrapeOK := true
	var err error

	e.scrapeMutex.Lock() // Do mutex unlock ASAP

	if e.scrapeActive {
		scrapeOK = false
		log.Warning("Collect is still active... - Skipping now")
		e.scrapeMutex.Unlock()
	} else {
		log.Debug("Collect started")

		e.scrapeActive = true
		e.scrapeMutex.Unlock()

		var start time.Time
		var elapsed float64

		e.stageExecutionMetric.Reset()
		e.jobMetadataOperationsMetric.Reset()
		e.jobReadThroughputMetric.Reset()
		e.jobWriteThroughputMetric.Reset()
		e.procMetadataOperationsMetric.Reset()
		e.procReadThroughputMetric.Reset()
		e.procWriteThroughputMetric.Reset()

		go retrieveRunningJobs(e.channelRunningJobs)
		go createUserInfoMap(e.channelUserInfo)
		go createGroupInfoMap(e.channelGroupInfo)

		runningJobsResult := <-e.channelRunningJobs
		userInfoResult := <-e.channelUserInfo
		groupInfoResult := <-e.channelGroupInfo

		if runningJobsResult.err != nil {
			scrapeOK = false
			log.Errorln(runningJobsResult.err)
		}
		if userInfoResult.err != nil {
			scrapeOK = false
			log.Errorln(userInfoResult.err)
		}
		if groupInfoResult.err != nil {
			scrapeOK = false
			log.Errorln(groupInfoResult.err)
		}

		e.stageExecutionMetric.WithLabelValues("retrieve_running_jobs").Set(runningJobsResult.elapsed)
		e.stageExecutionMetric.WithLabelValues("retrieve_user_name_info").Set(userInfoResult.elapsed)
		e.stageExecutionMetric.WithLabelValues("retrieve_group_name_info").Set(groupInfoResult.elapsed)

		start = time.Now()
		err = e.buildLustreMetadataMetrics(runningJobsResult.jobs, userInfoResult.users, groupInfoResult.groups)

		elapsed = time.Since(start).Seconds()
		e.stageExecutionMetric.WithLabelValues("build_metadata_metrics").Set(elapsed)

		if err != nil {
			if scrapeOK {
				scrapeOK = false
			}
			log.Errorln(err)
		}

		start = time.Now()
		err = e.buildLustreThroughputMetrics(runningJobsResult.jobs, userInfoResult.users, groupInfoResult.groups, true)
		elapsed = time.Since(start).Seconds()
		e.stageExecutionMetric.WithLabelValues("build_read_throughput_metrics").Set(elapsed)

		if err != nil {
			if scrapeOK {
				scrapeOK = false
			}
			log.Errorln(err)
		}

		start = time.Now()
		err = e.buildLustreThroughputMetrics(runningJobsResult.jobs, userInfoResult.users, groupInfoResult.groups, false)
		elapsed = time.Since(start).Seconds()
		e.stageExecutionMetric.WithLabelValues("build_write_throughput_metrics").Set(elapsed)

		if err != nil {
			if scrapeOK {
				scrapeOK = false
			}
			log.Errorln(err)
		}

		e.stageExecutionMetric.Collect(ch)
		e.jobMetadataOperationsMetric.Collect(ch)
		e.jobReadThroughputMetric.Collect(ch)
		e.jobWriteThroughputMetric.Collect(ch)
		e.procMetadataOperationsMetric.Collect(ch)
		e.procReadThroughputMetric.Collect(ch)
		e.procWriteThroughputMetric.Collect(ch)

		e.scrapeActive = false

		log.Debug("Collect finished")
	}

	if scrapeOK {
		e.scrapeOKMetric.Set(1)
	} else {
		e.scrapeOKMetric.Set(0)
	}

	e.scrapeOKMetric.Collect(ch)
}

func (e *exporter) Describe(ch chan<- *prometheus.Desc) {
	e.scrapeOKMetric.Describe(ch)
	e.stageExecutionMetric.Describe(ch)
	e.jobMetadataOperationsMetric.Describe(ch)
	e.jobReadThroughputMetric.Describe(ch)
	e.jobWriteThroughputMetric.Describe(ch)
	e.procMetadataOperationsMetric.Describe(ch)
	e.procReadThroughputMetric.Describe(ch)
	e.procWriteThroughputMetric.Describe(ch)
}

func (e *exporter) buildLustreMetadataMetrics(jobs []jobInfo, users userInfoMap, groups groupInfoMap) error {

	log.Debug("Process metadata operations")

	if jobs == nil || len(jobs) == 0 {
		return errors.New("Parameter jobs is not set")
	}

	if users == nil || len(users) == 0 {
		return errors.New("Parameter users is not set")
	}

	if groups == nil || len(groups) == 0 {
		return errors.New("Parameter groups is not set")
	}

	content, err := httpRequest(e.urlLustreMetadataOperations, e.requestTimeout)
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Trace("Bytes received: ", len(*content))
	}

	lustreMetadataOperations := parseLustreMetadataOperations(content)

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Count Lustre Jobids with metadata operatons: ", len(*lustreMetadataOperations))
	}

	for _, metadataInfo := range *lustreMetadataOperations {

		if isNumber(&metadataInfo.jobid) { // SLURM Job

			for _, job := range jobs {
				if metadataInfo.jobid == job.jobid {
					e.jobMetadataOperationsMetric.WithLabelValues(job.account, job.user).Add(
						float64(metadataInfo.operations))
				}
			}

		} else { // Process with UID (proc_name.uid)

			fields := strings.Split(metadataInfo.jobid, ".")
			lenFields := len(fields)

			var procName, userName, groupName string
			var uid int

			if lenFields == 2 {
				procName = fields[0]

				uid, err = strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
			} else if lenFields > 2 {
				lastFieldIdx := lenFields - 1
				procName = strings.Join((fields[0:lastFieldIdx]), ".")

				uid, err = strconv.Atoi(fields[lastFieldIdx])
				if err != nil {
					return err
				}
			} else {
				return errors.New("To few Lustre Jobstats procname_uid fields: " + metadataInfo.jobid)
			}

			userInfo, ok := users[uid]
			if !ok {
				return errors.New("uid not found in users map: " + strconv.Itoa(uid))
			}

			groupInfo, ok := groups[userInfo.gid]
			if !ok {
				return errors.New("gid not found in groups map: " + strconv.Itoa(userInfo.gid))
			}

			userName = userInfo.user
			groupName = groupInfo.group
			e.procMetadataOperationsMetric.WithLabelValues(procName, groupName, userName).Add(
				float64(metadataInfo.operations))
		}
	}

	return nil
}

func (e *exporter) buildLustreThroughputMetrics(jobs []jobInfo, users userInfoMap, groups groupInfoMap, read bool) error {

	var url string
	var jobMetric *prometheus.GaugeVec
	var procMetric *prometheus.GaugeVec

	if read {
		log.Debug("Process read throughput")
		url = e.urlLustreJobReadBytes
		jobMetric = e.jobReadThroughputMetric
		procMetric = e.procReadThroughputMetric
	} else {
		log.Debug("Process write throughput")
		url = e.urlLustreJobWriteBytes
		jobMetric = e.jobWriteThroughputMetric
		procMetric = e.procWriteThroughputMetric
	}

	if jobs == nil || len(jobs) == 0 {
		return errors.New("Parameter jobs is not set")
	}

	if users == nil || len(users) == 0 {
		return errors.New("Parameter users is not set")
	}

	if groups == nil || len(groups) == 0 {
		return errors.New("Parameter groups is not set")
	}

	content, err := httpRequest(url, e.requestTimeout)
	if err != nil {
		return err
	}

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Trace("Bytes received: ", len(*content))
	}

	lustreThroughput := parseLustreTotalBytes(content)

	if log.IsLevelEnabled(log.DebugLevel) {
		log.Debug("Count Lustre Jobids with throughput: ", len(*lustreThroughput))
	}

	for _, thInfo := range *lustreThroughput {

		if isNumber(&thInfo.jobid) { // SLURM Job

			for _, job := range jobs {
				if thInfo.jobid == job.jobid {
					jobMetric.WithLabelValues(job.account, job.user).Add(thInfo.throughput)
				}
			}

		} else { // Process with UID (proc_name.uid)

			fields := strings.Split(thInfo.jobid, ".")
			lenFields := len(fields)

			var procName, userName, groupName string
			var uid int

			if lenFields == 2 {
				procName = fields[0]

				uid, err = strconv.Atoi(fields[1])
				if err != nil {
					return err
				}
			} else if lenFields > 2 {
				lastFieldIdx := lenFields - 1
				procName = strings.Join((fields[0:lastFieldIdx]), ".")

				uid, err = strconv.Atoi(fields[lastFieldIdx])
				if err != nil {
					return err
				}
			} else {
				return errors.New("To few Lustre Jobstats procname_uid fields: " + thInfo.jobid)
			}

			userInfo, ok := users[uid]
			if !ok {
				return errors.New("uid not found in users map: " + strconv.Itoa(uid))
			}

			groupInfo, ok := groups[userInfo.gid]
			if !ok {
				return errors.New("gid not found in groups map: " + strconv.Itoa(userInfo.gid))
			}

			userName = userInfo.user
			groupName = groupInfo.group
			procMetric.WithLabelValues(procName, groupName, userName).Add(thInfo.throughput)
		}
	}

	return nil
}

func parseLustreMetadataOperations(content *[]byte) *[]metadataInfo {

	log.Debug("Parsing Lustre metadata operations")

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Trace(string(*content))
	}

	status, err := jsonparser.GetString(*content, "status")
	if err != nil || status != "success" {
		log.Panic(err)
	}

	slice := make([]metadataInfo, 0, 1000)

	jsonparser.ArrayEach(*content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		jobid, err := jsonparser.GetString(value, "metric", "jobid")

		if err != nil {
			// Might be the case with the exported Lustre jobstats. Cause not clear, need to check Lustre exporter.
			log.Warning("Key jobid not found in metric value:", string(value))
		} else {
			operationsStr, err := jsonparser.GetString(value, "value", "[1]")
			if err != nil {
				log.Panic(err)
			}

			operations, err := strconv.ParseInt(operationsStr, 10, 64)
			if err != nil {
				log.Panic(err)
			}

			slice = append(slice, metadataInfo{jobid, operations})
		}

	}, "data", "result")

	return &slice
}

func parseLustreTotalBytes(content *[]byte) *[]throughputInfo {

	log.Debug("Parsing Lustre total bytes")

	if log.IsLevelEnabled(log.TraceLevel) {
		log.Trace(string(*content))
	}

	status, err := jsonparser.GetString(*content, "status")
	if err != nil || status != "success" {
		log.Panic(err)
	}

	slice := make([]throughputInfo, 0, 1000)

	jsonparser.ArrayEach(*content, func(value []byte, dataType jsonparser.ValueType, offset int, err error) {

		jobid, err := jsonparser.GetString(value, "metric", "jobid")

		if err != nil {
			// Might be the case with the exported Lustre jobstats. Cause not clear, need to check Lustre exporter.
			log.Warning("Key jobid not found in metric value:", string(value))
		} else {
			throughputStr, err := jsonparser.GetString(value, "value", "[1]")
			if err != nil {
				log.Panic(err)
			}

			throughput, err := strconv.ParseFloat(throughputStr, 64)
			if err != nil {
				log.Panic(err)
			}
			slice = append(slice, throughputInfo{jobid, throughput})
		}

	}, "data", "result")

	return &slice
}

func isNumber(input *string) bool {
	if _, err := strconv.Atoi(*input); err != nil {
		return false
	}
	return true
}
