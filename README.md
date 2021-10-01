# Prometheus Cluster Exporter

A [Prometheus](https://prometheus.io/) exporter for Lustre metadata operations and IO throughput metrics associated to SLURM accounts  
and process names with user and group information on a cluster.

[Grafana dashboard](https://grafana.com/grafana/dashboards/14668) is also available.

## Getting

`go get github.com/GSI-HPC/prometheus-cluster-exporter`

## Building

```
cd $GOPATH/src/github.com/GSI-HPC/prometheus-cluster-exporter
go build -o prometheus-cluster-exporter *.go
```

## Requirements

### Lustre Exporter

[Lustre exporter](https://github.com/GSI-HPC/lustre_exporter) that exposes enabled Lustre Jobstats on the filesystem.

### Squeue Command

The squeue command from SLURM must be accessable locally to the exporter to retrieve the running jobs.  

For instance running the exporter on the SLURM controller is advisable, since the target host should be  
most stable for a productional environment.

### Getent

The getent command is required for the uid to user and group mapping used for the process names throughput metrics.

## Execution

### Parameter

| Name       | Default           | Description                                                                                                                        |
| ---------- | ----------------- | ---------------------------------------------------------------------------------------------------------------------------------- |
| version    | false             | Print version                                                                                                                      | 
| promserver |                   | [REQUIRED] Prometheus Server to be used e.g. http://prometheus-server:9090                                                         |
| log        | INFO              | Sets log level - INFO, DEBUG or TRACE                                                                                              | 
| port       | 9846              | The port to listen on for HTTP requests                                                                                            |
| timeout    | 15                | HTTP request timeout in seconds for exporting Lustre Jobstats on Prometheus HTTP API                                               |
| timerange  | 1m                | Time range used for rate function on the retrieving Lustre metrics from Prometheus - A three digit number with unit s, m, h or d   |

### Running in a Productive Environment

For a productive environment it is advisable to run the exporter on the SLURM controller,  
since the target host should be most stable.

### Prometheus Scrape Settings

Depending on the required resolution and runtime of the exporter,  
* the `scrape interval` should be set as appropriate e.g. at least 1 minute or higher.  
* the `scrape timeout` should be set close to the specified scrape interval.

## Metrics

Cluster exporter metrics are prefixed with "cluster_".

### Global

These metrics are always exported.

| Metric                              | Labels        | Description                                                       |
| ----------------------------------- | ------------- | ----------------------------------------------------------------- |
| exporter\_scrape\_ok                | -             | Indicates if the scrape of the exporter was successful or not.    |
| exporter\_stage\_execution\_seconds | name          | Execution duration in seconds spend in a specific exporter stage. |

### Metadata

#### **Jobs**

| Metric                     | Labels        | Description                                                 |
| ---------------------------| ------------- | ----------------------------------------------------------- |
| job\_metadata\_operations  | account, user | Total metadata operations of all jobs per account and user. |

#### **Process Names**

| Metric                     | Labels                              | Description                                                     |
| -------------------------- | ----------------------------------- | --------------------------------------------------------------- |
| proc\_metadata\_operations | proc\_name, group\_name, user\_name | Total metadata operations of process names per group and user.  |


### Throughput

#### **Jobs**

| Metric                        | Labels        | Description                                                                           |
| ----------------------------- | ------------- | ------------------------------------------------------------------------------------- |
| job\_read\_throughput\_bytes  | account, user | Total IO read throughput of all jobs on the cluster per account in bytes per second.  |
| job\_write\_throughput\_bytes | account, user | Total IO write throughput of all jobs on the cluster per account in bytes per second. |

#### **Process Names**

| Metric                         | Labels                              | Description                                                                                       |
| ------------------------------ | ----------------------------------- | ------------------------------------------------------------------------------------------------- |
| proc\_read\_throughput\_bytes  | proc\_name, group\_name, user\_name | Total IO read throughput of process names on the cluster per group and user in bytes per second.  |
| proc\_write\_throughput\_bytes | proc\_name, group\_name, user\_name | Total IO write throughput of process names on the cluster per group and user in bytes per second. |

## Multiple Srape Prevention

Since the forked processes do not have a timeout handling, they might block for a uncertain amount of time.  
It is very unlikely that reexecuting the processes will solve the problem of beeing blocked.  
Therefore multiple scrapes at a time will be prevented by the exporter.  

The following warning will be displayed on afterward scrape executions, were a scrape is still active:  
    *"Collect is still active... - Skipping now"*

Besides that, the cluster\_exporter\_scrape\_ok metric will be set to 0 for skipped scrape attempts.  

