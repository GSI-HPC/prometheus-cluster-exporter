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
	"flag"
	"net/http"
	"os"
	"regexp"
	"strconv"
	"strings"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

const (
	version                 = "1.1.4"
	namespace               = "cluster"
	namespaceInternals      = "cluster_exporter"
	httpApi                 = "/api/v1/query"
	queryParameter          = "?query=" // Query parameter are encoded in hex with % and 2 digits in the URL.
	queryMetadataOperations = "round%28sum%20by%28jobid%29%28irate%28lustre_job_stats_total[__TIME_RANGE__]%29%3E=1%29%29"
	queryJobReadBytes       = "sum%20by%28jobid%29%28irate%28lustre_job_read_bytes_total[__TIME_RANGE__]%29!=0%29"
	queryJobWriteBytes      = "sum%20by%28jobid%29%28irate%28lustre_job_write_bytes_total[__TIME_RANGE__]%29!=0%29"
	defaultLogLevel         = "INFO"
	defaultPort             = "9846"
	defaultRequestTimeout   = 15
	defaultTimeRange        = "1m"
)

type urlExportLustreMetrics struct {
	metadataOperations string
	jobReadBytes       string
	jobWriteBytes      string
}

func initLogging(logLevel string) {

	if logLevel == "INFO" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "TRACE" {
		log.SetLevel(log.TraceLevel)
	} else {
		log.Panicln("Not supported log level set")
	}

	log.SetOutput(os.Stdout)
}

func validateTimeRange(timeRange string) {

	lenTimeRange := len(timeRange)

	if lenTimeRange < 2 || lenTimeRange > 4 {
		log.Panic("Time range length is not supported: ", timeRange)
	}

	reTimeRangeUnit := regexp.MustCompile("s|m|h|d")

	timeRangeUnit := timeRange[lenTimeRange-1:]
	timeRangeNumber := timeRange[:lenTimeRange-1]

	if !reTimeRangeUnit.MatchString(timeRangeUnit) {
		log.Panic("Time range unit is not supported: ", timeRangeUnit)
	}

	_, err := strconv.Atoi(timeRangeNumber)

	if err != nil {
		log.Panic("Time range number could not be coverted to an integer: ", timeRangeNumber)
	}
}

func newUrlExportLustreMetrics(server string, timeRange string) *urlExportLustreMetrics {

	validateTimeRange(timeRange)

	serverQueryEndpoint := server + httpApi + queryParameter

	return &urlExportLustreMetrics{
		metadataOperations: serverQueryEndpoint + strings.Replace(queryMetadataOperations, "__TIME_RANGE__", timeRange, 1),
		jobReadBytes:       serverQueryEndpoint + strings.Replace(queryJobReadBytes, "__TIME_RANGE__", timeRange, 1),
		jobWriteBytes:      serverQueryEndpoint + strings.Replace(queryJobWriteBytes, "__TIME_RANGE__", timeRange, 1),
	}
}

func main() {

	printVersion := flag.Bool("version", false, "Print version")
	promServer := flag.String("promserver", "", "[REQUIRED] Prometheus Server to be used e.g. http://prometheus-server:9090")
	logLevel := flag.String("log", defaultLogLevel, "Sets log level - INFO, DEBUG or TRACE")
	port := flag.String("port", defaultPort, "The port to listen on for HTTP requests")
	requestTimeout := flag.Int("timeout", defaultRequestTimeout, "HTTP request timeout in seconds for exporting Lustre Jobstats on Prometheus HTTP API")
	timeRange := flag.String("timerange", defaultTimeRange, "Time range used for rate function on the retrieving Lustre metrics from Prometheus - A three digit number with unit s, m, h or d")

	flag.Parse()

	initLogging(*logLevel)

	if *printVersion {
		log.Info("Version: ", version)
		os.Exit(0)
	}

	if *promServer == "" {
		log.Panic("No Prometheus server has been specified")
	}

	metricsPath := "/metrics"
	listenAddress := ":" + *port

	log.Info("Exporter started")

	urlExports := newUrlExportLustreMetrics(*promServer, *timeRange)

	e := newExporter(*requestTimeout, urlExports.metadataOperations, urlExports.jobReadBytes, urlExports.jobWriteBytes)
	prometheus.MustRegister(e)

	http.Handle(metricsPath, promhttp.Handler())

	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Write([]byte(`<html>
             <head><title>Cluster Exporter</title></head>
             <body>
             <h1>Cluster Exporter</h1>
             <p><a href='` + metricsPath + `'>Metrics</a></p>
             </body>
             </html>`))
	})

	if err := http.ListenAndServe(listenAddress, nil); err != nil {
		log.Error(err)
	}

	log.Info("Exporter finished")
}
