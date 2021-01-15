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
	"flag"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	log "github.com/sirupsen/logrus"
)

const (
	version                       = "1.1"
	namespace                     = "cluster"
	namespaceInternals            = "cluster_exporter"
	defaultPort                   = "9166" // Port already in use by Dovecot exporter (https://github.com/prometheus/prometheus/wiki/Default-port-allocations)
	defaultRequestTimeout         = 15
	defaultLogLevel               = "INFO" // [INFO || DEBUG || TRACE]
	defaultURLLustreJobReadBytes  = "http://lustre-monitoring.gsi.de:9090/api/v1/query?query=sum%20by%28jobid%29%28irate%28lustre_job_read_bytes_total[1m]%29!=0%29"
	defaultURLLustreJobWriteBytes = "http://lustre-monitoring.gsi.de:9090/api/v1/query?query=sum%20by%28jobid%29%28irate%28lustre_job_write_bytes_total[1m]%29!=0%29"
)

func initLogging(logLevel string) {

	if logLevel == "INFO" {
		log.SetLevel(log.InfoLevel)
	} else if logLevel == "DEBUG" {
		log.SetLevel(log.DebugLevel)
	} else if logLevel == "TRACE" {
		log.SetLevel(log.TraceLevel)
	} else {
		log.Panicln("Not supported log level set!")
	}

	log.SetOutput(os.Stdout)
}

func main() {

	printVersion := flag.Bool("version", false, "Print version")
	logLevel := flag.String("log", defaultLogLevel, "Sets log level - INFO, DEBUG or TRACE")
	port := flag.String("port", defaultPort, "The port to listen on for HTTP requests.")
	requestTimeout := flag.Int("timeout", defaultRequestTimeout, "HTTP request timeout for exporting Lustre Jobstats on Prometheus HTTP API")
	urlLustreJobReadBytes := flag.String("urlReads", defaultURLLustreJobReadBytes, "URL with the query to the Prometheus HTTP API that exports the aggregated Lustre jobstats for the lustre_job_read_bytes_total")
	urlLustreJobWriteBytes := flag.String("urlWrites", defaultURLLustreJobWriteBytes, "URL with the query to the Prometheus HTTP API that exports the aggregated Lustre jobstats for the lustre_job_write_bytes_total")

	flag.Parse()

	if *printVersion {
		log.Info("Version: ", version)
		os.Exit(0)
	}

	initLogging(*logLevel)

	metricsPath := "/metrics"
	listenAddress := ":" + *port

	log.Info("Exporter started")

	e := newExporter(*requestTimeout, *urlLustreJobReadBytes, *urlLustreJobWriteBytes)
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
