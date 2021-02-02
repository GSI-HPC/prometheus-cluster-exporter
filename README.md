# Prometheus Cluster Exporter

A Prometheus exporter for Lustre IO throughput metrics associated to SLURM accounts and process names on a cluster.

## Requirements

### Lustre Exporter

A Lustre exporter that exposes the two metrics to Prometheus with a label jobid is required:

* lustre\_job\_read\_bytes\_total
* lustre\_job\_write\_bytes\_total

The Lustre exporter from HP provides such metrics:
https://github.com/HewlettPackard/lustre\_exporter

### Squeue Command

The squeue command from SLURM must be accessable locally to the exporter to retrieve the running jobs.

### Getent

The getent command is required for the uid to user and group mapping used for the process names throughput metrics.

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

## Building the Exporter

```go
go build -o cluster-exporter *.go
```
