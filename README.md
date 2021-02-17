# Prometheus Cluster Exporter

A Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts  
and process names with user and group information on a cluster.

## Building

`go build -o cluster_exporter *.go`

## Requirements

### Lustre Exporter

A Lustre exporter that exposes the two metrics to Prometheus with a label jobid is required:

* lustre\_job\_read\_bytes\_total
* lustre\_job\_write\_bytes\_total

The Lustre exporter from HP provides such metrics: https://github.com/HewlettPackard/lustre\_exporter

### Squeue Command

The squeue command from SLURM must be accessable locally to the exporter to retrieve the running jobs.  

For instance running the exporter on the SLURM controller is advisable, since the target host should be  
most stable for a productional environment.

### Getent

The getent command is required for the uid to user and group mapping used for the process names throughput metrics.

## Execution

### Parameter

| Name      | Default           | Description                                                                                                                        |
| --------- | ----------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| log       | INFO              | Logging level                                                                                                                      | 
| port      | 9166              | The port to listen on for HTTP requests                                                                                            |
| timeout   | 15                | HTTP request timeout in seconds for exporting Lustre Jobstats on Prometheus HTTP API                                               |
| urlReads  | Site specific URL | Query URL to the Prometheus HTTP API that exports the Lustre jobstats read throughput rate                                         |
| urlWrites | Site specific URL | Query URL to the Prometheus HTTP API that exports the Lustre jobstats write throughput rate                                        |

### Exporting Lustre Jobstats Throughput Rate

The Lustre jobstats throughput rates are calculated on the Prometheus server and exported via HTTP API.  

A HTTP query for setting the urlReads and urlWrites parameter consists of:

* Server endpoint: `http://prom-server:9090/`
* HTTP API endpoint: `api/v1/query`
* Query parameter: `?query=`
* Query string
    * Reads: `sum%20by%28jobid%29%28irate%28lustre_job_write_bytes_total[1m]%29!=0%29`
    * Writes: `sum%20by%28jobid%29%28irate%28lustre_job_write_bytes_total[1m]%29!=0%29`

> Some special character are defined in UTF-8 hexadecimal in the query strings.

### Running in a Productive Environment

For a productive environment it is advisable to run the exporter on the SLURM controller,  
since the target host should be most stable.

### Prometheus Scrape Settings

Depending on the required resolution and runtimes of an exporter  
the `scrape interval` should be set as appropriate e.g. at least 1m or higher.  

The `scrape timeout` should also be set, but close to the specific scrape interval.

## Metrics

Cluster exporter metrics are prefixed with "cluster_".

### Global

These metrics are always exported.

| Metric                              | Labels        | Description                                                       |
| ----------------------------------- | ------------- | ----------------------------------------------------------------- |
| exporter\_scrape\_ok                | -             | Indicates if the scrape of the exporter was successful or not.    |
| exporter\_stage\_execution\_seconds | name          | Execution duration in seconds spend in a specific exporter stage. |

### Throughput

#### **Jobs**

| Metric                        | Labels        | Description                                                                           |
| ----------------------------- | ------------- | ------------------------------------------------------------------------------------- |
| job\_read\_throughput\_bytes  | account, user | Total IO read throughput of all jobs on the cluster per account in bytes per second.  |
| job\_write\_throughput\_bytes | account, user | Total IO write throughput of all jobs on the cluster per account in bytes per second. |

#### **Process Names**

| Metric                         | Labels                              | Description                                                                                       |
| ------------------------------ | ----------------------------------- | ------------------------------------------------------------------------------------------------- |
| proc\_read\_throughput\_bytes  | proc\_name, user\_name, group\_name | Total IO read throughput of process names on the cluster per group and user in bytes per second.  |
| proc\_write\_throughput\_bytes | proc\_name, user\_name, group\_name | Total IO write throughput of process names on the cluster per group and user in bytes per second. |

## Multiple Srape Prevention

Since the forked processes do not have a timeout handling, they might block for a uncertain amount of time.  
It is very unlikely that reexecuting the processes will solve the problem of beeing blocked.  
Therefore multiple scrapes at a time will be prevented by the exporter.  

The following warning will be displayed on afterward scrape executions, were a scrape is still active:  
    *"Collect is still active... - Skipping now"*

Besides that, the cluster\_exporter\_scrape\_ok metric will be set to 0 for skipped scrape attempts.  

